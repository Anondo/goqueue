package queue

import (
	"github.com/spf13/cobra"
)

var (
	QueueCmd = &cobra.Command{
		Use:   "queue",
		Short: "Commands regrading queues",
	}
)

func init() {
	QueueCmd.AddCommand(ListCmd)
	QueueCmd.AddCommand(RmCmd)
}
