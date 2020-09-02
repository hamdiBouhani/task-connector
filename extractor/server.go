package extractor

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/target-smart-data-ai-searsh/task-connector-be/config"
	"gitlab.com/target-smart-data-ai-searsh/task-connector-be/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Server struct.
type Server struct {
	*logrus.Logger
	*gin.Engine

	WsAddr string //                         Url for  ws Qraphgl API

	Session *types.Session

	oidcAddr     string
	clientID     string
	clientSecret string

	DbName string
	PgUser string
	PgPass string
	PgHost string
	PgPort int

	mongoClt  *mongo.Client
	mongoDB   *mongo.Database
	mongoAddr string
}

//go run cmd/main.go extract-tasks --db=mws --service-account="./service-account-file.json"

//NewExtractor create new server struct.
func NewExtractor(
	WsAddr string,
	dbName string,
	pgUser string,
	pgPass string,
	pgHost string,
	pgPort int,
	topicName string, /*The cloud pubsub topic to use*/
	projectID string, /* Google Cloud Platform project ID.*/
	serviceAccount string, /*"service-account Path to Service Account .json file"*/
	clientID string,
	clientSecret string,
	oidcAddr string,
	mongoAddr string,
) (*Server, error) {

	//Logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	//session
	session, err := config.InitSession(dbName, pgUser, pgPass, pgHost, pgPort)
	if err != nil {
		e := fmt.Sprintf("unable to init session")
		logger.WithError(err).Error(e)
		return nil, errors.Wrap(err, e)
	}

	//pubsub
	session.PSConn = config.NewPubSubCon(topicName, projectID, serviceAccount)

	//Mongo
	MongoClient, err := mongo.NewClient(options.Client().ApplyURI(mongoAddr))

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	logger.Info("connect to Mongodb ")
	err = MongoClient.Connect(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to mongo db")
	}

	db := MongoClient.Database("wsTasks")
	if db == nil {
		return nil, errors.Wrap(err, "faild to create notifications db")
	}

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// TODO uncomment + pass the argCORSHosts to the Header instead of *
		//c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Mode, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	})

	pprof.Register(r)

	server := &Server{
		WsAddr: WsAddr,

		DbName: dbName,
		PgUser: pgUser,
		PgPass: pgPass,
		PgHost: pgHost,
		PgPort: pgPort,

		Logger: logger,
		Engine: r,

		Session: session,

		oidcAddr:     oidcAddr,
		clientID:     clientID,
		clientSecret: clientSecret,

		mongoClt:  MongoClient,
		mongoDB:   db,
		mongoAddr: mongoAddr,
	}

	// open accessed group
	openAccessed := r.Group("/")
	{
		// service info handler
		openAccessed.GET("/info", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"info": "UP",
			})
		})
	}

	v1 := r.Group("/v1/api")
	{
		ws := v1.Group("/ws")
		ws.GET("/Stream", server.StreamChanges)
		ws.GET("/Snapshot", server.Snapshot)
	}

	return server, nil
}
