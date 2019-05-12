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

// Queue represents the task job queue
type Queue struct {
	ID                  int
	Name                string
	Jobs                chan Job
	Capacity            int
	RegisteredTaskNames []string
	Subscribers         []*Subscriber
	Durable             bool
	AckWait             time.Duration
}

// QList is the list of all the queues in the  broker
var (
	QList []Queue
)

// Init initializes the queues, both persistant and default queues
func Init() {
	durableFileName = viper.GetString("persistance.filepath")
	InitDefaultQueue()
	initPersistedQueues()

}

// PushTask pushes a new task into the queue
func (q *Queue) PushTask(jn string, args []Arguments) {
	j := Job{
		ID:      len(q.Jobs) + 1,
		JobName: jn,
		Args:    args,
	}

	jl := len(q.Jobs) + 1
	q.Jobs <- j

	helper.FailOnError(q.addDurableJob(j), "Could not persist job")

	helper.JobReceiveLog(j.JobName, q.Name, jl, q.Capacity, j.Args, q.Durable)

}

// IsTaskRegistered determines if a task is registered against the queue
func (q *Queue) IsTaskRegistered(tn string) bool {
	for _, t := range q.RegisteredTaskNames {
		if tn == t {
			return true
		}
	}
	return false
}

// Clear empties a queue
func (q *Queue) Clear() string {

	if len(q.Jobs) == 0 {
		return fmt.Sprintf("Queue:%s already empty", q.Name)
	}

	for len(q.Jobs) > 0 {
		<-q.Jobs
	}
	helper.ColorLog("\033[35m", fmt.Sprintf("Queue Cleared: {Name:%s & Capacity:%d}\n", q.Name, q.Capacity))

	return fmt.Sprintf("Queue:%s has been successfully truncated", q.Name)
}

// InitDefaultQueue initializes the default queue
func InitDefaultQueue() {
	q := Queue{
		ID:       1,
		Name:     viper.GetString("default.queue_name"),
		Capacity: viper.GetInt("default.queue_capacity"),
	}
	q.Jobs = make(chan Job, q.Capacity)
	QList = append(QList, q)
}

// AddQueue creates a new queue in the broker
func AddQueue(q Queue) {
	if QueueExists(q.Name) {
		return
	}
	QList = append(QList, q)
	if q.Durable {
		helper.FailOnError(q.persistQueue(), "Could not persist queue:"+q.Name)
	}
	helper.ColorLog("\033[35m", fmt.Sprintf("Queue Declared: {Name:%s , Capacity:%d , Durable:%v , AckWait:%v}\n",
		q.Name, q.Capacity, q.Durable, q.AckWait))
	return
}

// AddTask makes the queue push a new task
func AddTask(qn, jn string, args []Arguments) {
	for _, q := range QList {
		if q.Name == qn {
			q.PushTask(jn, args)
			return
		}
	}
	log.Println("No queue named", qn)
}

// QueueExists determines if a queue actually exists by its name
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
				helper.FailOnError(q.addDurableRegTask(tn), "Failed to persist task")
				helper.ColorLog("\033[35m", fmt.Sprintf("Successfully registered task:%v", tn))
			}
		}

		return nil
	}
	return errors.New("No Such Queue")
}

func (q *Queue) GetSubscriber(s string) *Subscriber {
	for _, sbscrbr := range q.Subscribers {
		if sbscrbr.CName == s {
			return sbscrbr
		}
	}
	return nil
}
