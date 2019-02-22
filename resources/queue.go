package resources

import (
	"context"
	"fmt"
	"goqueue/helper"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Queue struct {
	ID                  int
	Name                string
	Jobs                chan Job
	Capacity            int
	RegisteredTaskNames []string
	Subscribers         []*Subscriber
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

	jl := len(q.Jobs) + 1
	q.Jobs <- j

	helper.JobReceiveLog(j.JobName, q.Name, jl, q.Capacity, j.Args)

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
	if QueueExists(q.Name) {
		return
	}
	QList = append(QList, q)
	helper.ColorLog("\033[35m", fmt.Sprintf("Queue Declared: {Name:%s & Capacity:%d}\n", q.Name, q.Capacity))
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

func QueueExists(qn string) bool {
	for _, q := range QList {
		if q.Name == qn {
			return true
		}
	}
	return false
}

func GetQueueByName(qn string) *Queue {
	for i, q := range QList {
		if q.Name == qn {
			return &QList[i]
		}
	}
	return nil
}

func GetAck(q *Queue, hn, wn string) (bool, error) {
	for _, s := range q.Subscribers {
		if s.Host == hn {
			uri := "http://" + s.Host + ":" + strconv.Itoa(s.Port) + "/worker/acknowledge"
			req, err := http.NewRequest(http.MethodGet, uri, nil)

			if err != nil {
				return false, nil
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			req = req.WithContext(ctx)

			defer cancel()

			c := http.Client{}
			resp, _ := c.Do(req)

			if resp != nil {
				helper.ColorLog("\033[35m", fmt.Sprintf("Received acknowledged from consumer:%s", wn))
				return true, nil
			}

			break
		}
	}

	return false, nil
}
