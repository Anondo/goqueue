package resources

import (
	"fmt"
	"goqueue/helper"
	"log"
)

type Queue struct {
	ID       int
	Name     string
	Jobs     chan Job
	Capacity int
}

var (
	QList []Queue
)

func Init() {

	InitDefaultQueue()

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
	msg := fmt.Sprintf("Job Received: {Name: %s Args: %v}", j.JobName, j.Args)
	helper.ColorLog(prpl, msg)

	msg = fmt.Sprintf("Queue: %s", q.Name)
	helper.ColorLog(prpl, msg)

	msg = fmt.Sprintf("Total Jobs: %d", len(q.Jobs))
	helper.ColorLog(prpl, msg)

	msg = fmt.Sprintf("Queue Capacity: %d", q.Capacity)
	helper.ColorLog(prpl, msg)

	fmt.Println("--------------------------------------------")

}

func InitDefaultQueue() {
	q := Queue{
		ID:       1,
		Name:     "default_queue",
		Capacity: 1000,
	}
	q.Jobs = make(chan Job, q.Capacity)
	QList = append(QList, q)
}

func AddQueue(q Queue) {
	QList = append(QList, q)
}

func AddTask(qn, jn string, args []interface{}) {
	for _, q := range QList {
		if q.Name == qn {
			q.PushTask(jn, args)
			return
		}
	}
	log.Println("No queue named", qn)
}
