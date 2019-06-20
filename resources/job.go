package resources

import (
	"encoding/json"
	"fmt"
	"goqueue/helper"
	"net/http"
	"time"
)

// Arguments is the struct which will be used for taking arguments for a task
type Arguments struct {
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

// Job is the struct which represents a JOB queued in the queue
type Job struct {
	ID      string      `json:"id"`
	JobName string      `json:"job_name"`
	Args    []Arguments `json:"args"`
}

// SendJob is a function which sends the job to the subscribed consumer
func SendJob(w http.ResponseWriter, qn, wn string) {
	q := GetQueueByName(qn)
	if q != nil {

		j := <-q.Jobs

		s := q.GetSubscriber(wn)

		if s == nil { // when worker dies after making request,
			q.Jobs <- j // TODO: Need a better solution, for requeing for such situation
			return
		}

		if !q.IsTaskRegistered(j.JobName) { // TODO: Tasks are now registered regardless of the subscriber,need to fix this
			if q.Durable {
				helper.LogOnError(q.Requeue(), "Could not requeue task")
			}
			return
		}
		b, err := json.Marshal(j)
		helper.FailOnError(err, "Could not decode job")
		fmt.Fprintf(w, "%s", string(b))

		helper.ColorLog("\033[35m", fmt.Sprintf("Job:{Name:%s Args:%v} fetched by %s", j.JobName, j.Args, wn))

		go func() {
			ackEndTime := time.Now().Add(q.AckWait)

			for time.Now().Before(ackEndTime) {
				if s.Ack {
					helper.ColorLog("\033[35m", fmt.Sprintf("Received acknowledgement from consumer:%s", wn))
					helper.FailOnError(q.removeDurableJob(j), "Could not remove persistant job")
					helper.ColorLog("\033[35m", fmt.Sprintf("Subscriber:%s is ready to fetch jobs", wn))
					break
				}
			}

			if !s.Ack && q.Durable {
				helper.LogOnError(q.Requeue(), "Could not requeue task")
			}

			if !s.Ack {
				helper.ColorLog("\033[31m", fmt.Sprintf("No acknowledgement received from consumer:%s", wn))
			}

			s.Ack = false
		}()

	}
}
