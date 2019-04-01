package cmd

import (
	"goqueue/cmd/queue"
	"goqueue/config"
	"goqueue/helper"
	"goqueue/resources"

	"github.com/spf13/cobra"
)

var (
	//RootCmd is the root command for the goqueue app
	RootCmd = &cobra.Command{
		Use:   "goqueue",
		Short: "Goqueue is a task/job queue",
		Long:  helper.Logo,
	}
)

func init() {
	config.Init()
	resources.Init()
	RootCmd.AddCommand(ServerCmd)
	RootCmd.AddCommand(queue.QueueCmd)
}

// Execute executes the root command for the app
func Execute() {
	helper.FailOnError(RootCmd.Execute(), "Failed to run goqueue")
}
