package cmd

import (
	"goqueue/cmd/queue"
	"goqueue/helper"
	"goqueue/resources"

	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "goqueue",
		Short: "Goqueue is a task/job queue",
		Long:  helper.Logo,
	}
)

func init() {
	resources.Init()
	RootCmd.AddCommand(ServerCmd)
	RootCmd.AddCommand(queue.QueueCmd)
}

func Execute() {
	helper.FailOnError(RootCmd.Execute(), "Failed to run goqueue")
}
