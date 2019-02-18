package cmd

import (
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
}

func Execute() {
	helper.FailOnError(RootCmd.Execute(), "Failed to run goqueue")
}
