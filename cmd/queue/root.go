package queue

import (
	"github.com/spf13/cobra"
)

var (
	//QueueCmd is the root of every command related to queues
	QueueCmd = &cobra.Command{
		Use:   "queue",
		Short: "Commands regrading queues",
	}
)

func init() {
	QueueCmd.AddCommand(ListCmd)
	QueueCmd.AddCommand(RmCmd)
	QueueCmd.AddCommand(ClearCmd)
}

// TODO: All the queue commands makes a http request to the server. Need to fix this,
// no http calls ,for example: basic iteration on the Qlist.
