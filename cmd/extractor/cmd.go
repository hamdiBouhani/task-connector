package extractor

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/target-smart-data-ai-searsh/task-connector-be/extractor"
)

var Cmd = &cobra.Command{
	Use:   "extract-tasks",
	Short: "Start extract tasks service",
	Run:   run,
}

var (
	dbName string
	pgUser string
	pgPass string
	pgHost string
	pgPort int

	serverHost string
	serverPort string

	mongoAddr string

	wsAddr string

	topic          string
	projectID      string
	serviceAccount string

	clientID     string
	clientSecret string
	oidcAddr     string
)

func init() {
	Cmd.Flags().StringVar(&dbName, "db", "", "Name of the database to connect to")
	Cmd.Flags().StringVar(&pgUser, "user", "postgres", "Postgres user name")
	Cmd.Flags().StringVar(&pgPass, "password", "postgres", "Postgres password")
	Cmd.Flags().StringVar(&pgHost, "pgHost", "localhost", "Postgres server hostname")
	Cmd.Flags().IntVar(&pgPort, "pgPort", 5432, "Postgres server port")

	Cmd.Flags().StringVar(&serverHost, "serverHost", "0.0.0.0", "Host to listen on")
	Cmd.Flags().StringVar(&serverPort, "serverPort", "8080", "Port to listen on")

	Cmd.Flags().StringVar(&mongoAddr, "mongo-addr", "mongodb://localhost:27017", "mongo hostname")

	Cmd.Flags().StringVar(&wsAddr, "ws", "https://test.meeraspace.com/graphql", "workspace provider hostname")

	Cmd.Flags().StringVar(&topic, "topic", "taskTopic", "The PubSub topic to use")
	Cmd.Flags().StringVar(&projectID, "project-id", "smart-meera", "Sets your Google Cloud Platform project ID.")
	Cmd.Flags().StringVar(&serviceAccount, "service-account", "", "Path to Service Account .json file") //./service-account-file.json

	Cmd.Flags().StringVar(&clientID, "client-id", "nz-mdr", "auth2 client id")
	Cmd.Flags().StringVar(&clientSecret, "client-secret", "DzXZxyDObSpsnR7qLqQ4p1LEVoIiE49e", "auth2 client secret")
	Cmd.Flags().StringVar(&oidcAddr, "oidc", "https://sso.test.meeraspace.com/token", "oidc provider hostname")

}

func run(cmd *cobra.Command, args []string) {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	s, err := extractor.NewExtractor(
		wsAddr,
		dbName,
		pgUser,
		pgPass,
		pgHost,
		pgPort,
		topic,
		projectID,
		serviceAccount,
		clientID,
		clientSecret,
		oidcAddr,
		mongoAddr,
	)
	if err != nil {
		logger.Fatalln("couldn't create extractor server:", err)
	}

	logger.Infof("Starting server for database %s; serving at %s:%s", dbName, serverHost, serverPort)
	if err = s.Run(serverHost + ":" + serverPort); err != nil {
		logger.Fatalln("server error:", err)
	}
}
