package resources

import (
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

	helper.JobReceiveLog(j.JobName, q.Name, len(q.Jobs), q.Capacity, j.Args)

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
	for _, v := range QList {
		if q.Name == v.Name {
			v = q
			return
		}
	}
	QList = append(QList, q)
	return
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
