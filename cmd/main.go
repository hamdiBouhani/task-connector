package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"gitlab.com/target-smart-data-ai-searsh/task-connector-be/cmd/extractor"
	"gitlab.com/target-smart-data-ai-searsh/task-connector-be/cmd/tasks"
	"gitlab.com/target-smart-data-ai-searsh/task-connector-be/version"
)

var (
	rootCmd *cobra.Command
)

func init() {
	rootCmd = &cobra.Command{
		Use:   "tag-connector",
		Short: "utilies and services",
		Long:  "Top level command for utilities and services of the task-connector-be app",
	}

	rootCmd.AddCommand(
		versionCmd,
		tasks.Cmd,
		extractor.Cmd,
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version and exit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(
			"tag-connector-be Version: %s \n API Version: %s \n Go Version: %s \n Go OS/ARCH: %s %s",
			version.Version,
			version.APIVersion,
			runtime.Version(),
			runtime.GOOS,
			runtime.GOARCH,
		)
	},
}
