package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"gitlab.com/target-smart-data-ai-searsh/task-connector-be/model"

	"cloud.google.com/go/pubsub"
	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type TokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type TaskJobRunner struct {
	Logger *logrus.Logger
	WsAddr string //                         Url for  ws Qraphgl API

	MongoClient *mongo.Client
	MongoDB     *mongo.Database
	MongoAddr   string

	oidcAddr     string
	clientID     string
	clientSecret string

	topicName string
	ProjectID string
	topic     *pubsub.Topic
	Client    *pubsub.Client
}

func NewTaskJobRunner(
	WsAddr string,
	MongoAddr string,
	projectID string,
	serviceAccount string,
	topicName string,
	clientID string,
	clientSecret string,
	oidcAddr string,
) (*TaskJobRunner, error) {
	// Logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	//Mongo
	/*MongoClient, err := mongo.NewClient(options.Client().ApplyURI(MongoAddr))

	logger.Info("connect to Mongodb ")
	err = MongoClient.Connect(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to mongo db")
	}

	db := MongoClient.Database("workspaces")
	if db == nil {
		return nil, errors.Wrap(err, "faild to create notifications db")
	}

	logger.Info("Ping to Mongodb ")
	err = MongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, errors.Wrap(err, "unable to ping to mongo db")
	}*/

	//Cloud PubSub

	c := context.Background()

	if projectID == "" {
		log.Fatalf("Failed to create client: %v", errors.New("project id should not be empty"))
	}

	var client *pubsub.Client
	if len(serviceAccount) > 0 {
		var aud string = "https://pubsub.googleapis.com/google.pubsub.v1.Publisher"

		keyBytes, err := ioutil.ReadFile(serviceAccount)
		if err != nil {
			log.Fatalf("Unable to read service account key file  %v", err)
		}

		tokenSource, err := google.JWTAccessTokenSourceFromJSON(keyBytes, aud)
		if err != nil {
			log.Fatalf("Error building JWT access token source: %v", err)
		}

		client, err = pubsub.NewClient(ctx, projectID, option.WithTokenSource(tokenSource))
		if err != nil {
			log.Fatalf("Could not create pubsub Client: %v", err)
		}

	} else {
		logger.Println("Path to Service Account .json file is empty")
		// Creates a client.
		var err error
		client, err = pubsub.NewClient(c, projectID)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
	}

	if topicName == "" {
		log.Fatalf("Sets the name for the new topic.: %v", errors.New("topic name should not be empty"))
	}

	// Creates the new topic.
	var topic *pubsub.Topic
	topic = client.Topic(topicName)

	logger.Printf("Topic %v created.\n", topic)

	return &TaskJobRunner{
		WsAddr: WsAddr,
		// MongoAddr:   MongoAddr,
		// MongoClient: MongoClient,
		// MongoDB:     db,

		Logger: logger,

		topic:        topic,
		topicName:    topicName,
		ProjectID:    projectID,
		Client:       client,
		oidcAddr:     oidcAddr,
		clientID:     clientID,
		clientSecret: clientSecret,
	}, nil
}

//Run implement of get all workspace api call.
func (r *TaskJobRunner) Run() {

	r.Logger.Info(fmt.Sprintf("\n\n\n\n\n Task Job Running at %v . \n\n\n\n\n", time.Now()))

	var tokenResp TokenResp

	payload := strings.NewReader("grant_type=client_credentials&client_id=" + r.clientID + "&client_secret=" + r.clientSecret + "&scope=openid email groups profile offline_access")

	reqToken, _ := http.NewRequest("POST", r.oidcAddr, payload)
	reqToken.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(reqToken)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &tokenResp)

	client := graphql.NewClient(r.WsAddr)

	// make a request
	req := graphql.NewRequest(ALLWorkspaceQuery)
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)

	// run it and capture the response
	var respData model.AllWorkspaces
	if err := client.Run(context.Background(), req, &respData); err != nil {
		log.Fatal(err)
	}

	file, _ := json.MarshalIndent(respData, "", " ")

	_ = ioutil.WriteFile("test.json", file, 0644)

	r.SendToPubSub(respData.Data)
}
