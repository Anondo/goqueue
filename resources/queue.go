package resources

import (
	"fmt"
	"goqueue/helper"
)

type Queue struct {
	ID       int
	Name     string
	Jobs     chan Job
	Capacity int
}

var (
	Q Queue
)

func Init() {
	Q.Name = "default_queue"
	Q.Capacity = 100
	Q.Jobs = make(chan Job, Q.Capacity)
}

func (q *Queue) PushTask(jn string, args []interface{}) {
	j := Job{
		ID:      len(q.Jobs) + 1,
		JobName: jn,
		Args:    args,
	}
	// helper.FailOnError(j.ProcessJob(f, args), "Failed to push task")

	q.Jobs <- j

	prpl := "\033[35m" // purple
	msg := fmt.Sprintf("Job Received: {Name: %s Args: %v}\nTotal Jobs: %d", j.JobName, j.Args, len(q.Jobs))

	txt := fmt.Sprintf(`{
				%s
		}`, msg)
	helper.ColorLog(prpl, txt)

}
