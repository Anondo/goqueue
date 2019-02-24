package resources

import (
	"context"
	"errors"
	"fmt"
	"goqueue/helper"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/viper"
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

func (q *Queue) PushTask(jn string, args []Arguments) {
	j := Job{
		ID:      len(q.Jobs) + 1,
		JobName: jn,
		Args:    args,
	}

	jl := len(q.Jobs) + 1
	q.Jobs <- j

	helper.JobReceiveLog(j.JobName, q.Name, jl, q.Capacity, j.Args)

}

func (q *Queue) IsTaskRegistered(tn string) bool {
	for _, t := range q.RegisteredTaskNames {
		if tn == t {
			return true
		}
	}
	return false
}

func InitDefaultQueue() {
	q := Queue{
		ID:       1,
		Name:     viper.GetString("default.queue_name"),
		Capacity: viper.GetInt("default.queue_capacity"),
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

func AddTask(qn, jn string, args []Arguments) {
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

			ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("requests.timeout")*time.Second)
			req = req.WithContext(ctx)

			defer cancel()

			c := http.Client{}
			resp, _ := c.Do(req)

			if resp != nil {
				helper.ColorLog("\033[35m", fmt.Sprintf("Received acknowledgement from consumer:%s", wn))
				return true, nil
			}

			break
		}
	}
	helper.ColorLog("\033[35m", fmt.Sprintf("No acknowledgement from consumer:%s", wn))
	return false, nil
}

func RegisterTasks(qn string, tns []string) error {
	q := GetQueueByName(qn)
	if q != nil {
		for _, tn := range tns {
			if !q.IsTaskRegistered(tn) {
				q.RegisteredTaskNames = append(q.RegisteredTaskNames, tn)
				helper.ColorLog("\033[35m", fmt.Sprintf("Successfully registered task:%v", tn))
			}
		}

		return nil
	}
	return errors.New("No Such Queue")
}
