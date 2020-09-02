package tasks

import (
	"github.com/bamzi/jobrunner"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gitlab.com/target-smart-data-ai-searsh/task-connector-be/tasks"
)

var Cmd = &cobra.Command{
	Use:   "start-job",
	Short: "Start A job runner for executing scheduled",
	Run:   run,
}

var (
	wsAddr         string
	mongoAddr      string
	topic          string
	projectID      string
	serviceAccount string
	clientID       string
	clientSecret   string
	oidcAddr       string
)

func init() {
	Cmd.Flags().StringVar(&wsAddr, "ws", "https://test.meeraspace.com/graphql", "workspace provider hostname")
	Cmd.Flags().StringVar(&mongoAddr, "mongo-addr", "mongodb://localhost:27017", "mongo hostname")
	Cmd.Flags().StringVar(&topic, "topic", "taskTopic", "The Kafka topic to use")
	Cmd.Flags().StringVar(&projectID, "project-id", "target-datalake-ng", "Sets your Google Cloud Platform project ID.")
	Cmd.Flags().StringVar(&serviceAccount, "service-account", "", "Path to Service Account .json file") //./service-account-file.json
	Cmd.Flags().StringVar(&clientID, "client-id", "nz-mdr", "auth2 client id")
	Cmd.Flags().StringVar(&clientSecret, "client-secret", "DzXZxyDObSpsnR7qLqQ4p1LEVoIiE49e", "auth2 client secret")
	Cmd.Flags().StringVar(&oidcAddr, "oidc", "https://sso.test.meeraspace.com/token", "oidc provider hostname")

}

func run(cmd *cobra.Command, args []string) {

	t, _ := tasks.NewTaskJobRunner(wsAddr, mongoAddr, projectID, serviceAccount, topic, clientID, clientSecret, oidcAddr)

	//t.Run()
	jobrunner.Start() // optional: jobrunner.Start(pool int, concurrent int) (10, 1)
	jobrunner.Schedule("@every 30m00s", t)

	routes := gin.Default()

	// Resource to return the JSON data
	routes.GET("/jobrunner/json", JobJson)
	routes.Run(":8080")
}

func JobJson(c *gin.Context) {
	// returns a map[string]interface{} that can be marshalled as JSON
	c.JSON(200, jobrunner.StatusJson())
}
