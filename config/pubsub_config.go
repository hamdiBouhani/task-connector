package config

import (
	"context"
	"errors"
	"io/ioutil"

	"cloud.google.com/go/pubsub"
	"github.com/sirupsen/logrus"
	"gitlab.com/target-smart-data-ai-searsh/task-connector-be/types"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

//NewPubSubCon create new PubSubCon struct.
func NewPubSubCon(topicName string, projectID string, serviceAccount string) *types.PubSubCon {
	var ConfigPubSub = types.PubSubCon{}
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)

	ctx := context.Background()

	if projectID == "" {
		log.Fatalf("Failed to create client: %v", errors.New("project id should not be empty"))
	}

	ConfigPubSub.ProjectID = projectID

	var client *pubsub.Client
	if len(serviceAccount) > 0 {

		aud := "https://pubsub.googleapis.com/google.pubsub.v1.Publisher"

		keyBytes, err := ioutil.ReadFile(serviceAccount)
		if err != nil {
			log.Fatalf("Unable to read service account key file  %v", err)
		}

		ConfigPubSub.ServiceAccount = serviceAccount

		tokenSource, err := google.JWTAccessTokenSourceFromJSON(keyBytes, aud)
		if err != nil {
			log.Fatalf("Error building JWT access token source: %v", err)
		}

		client, err = pubsub.NewClient(ctx, projectID, option.WithTokenSource(tokenSource))
		if err != nil {
			log.Fatalf("Could not create pubsub Client: %v", err)
		}

	} else {
		log.Println("Path to Service Account .json file is empty")

		// Creates a client.
		var err error
		client, err = pubsub.NewClient(ctx, projectID)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
	}

	ConfigPubSub.Client = client

	if topicName == "" {
		log.Fatalf("Sets the name for the new topic.: %v", errors.New("topic name should not be empty"))
	}

	ConfigPubSub.TopicName = topicName

	// Creates the new topic.
	var topic *pubsub.Topic
	topic = client.Topic(topicName)
	ConfigPubSub.Topic = topic

	log.Printf("Topic %v created.\n", topic)
	return &ConfigPubSub
}
