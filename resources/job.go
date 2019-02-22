package resources

import (
	"encoding/json"
	"fmt"
	"goqueue/helper"
	"net/http"
)

type Arguments struct {
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

type Job struct {
	ID      int         `json:"id"`
	JobName string      `json:"job_name"`
	Args    []Arguments `json:"args"`
}

func SendJob(w http.ResponseWriter, qn, wn, hn string) {
	q := GetQueueByName(qn)
	if q != nil {

		helper.ColorLog("\033[35m", fmt.Sprintf("Subscriber:%s is ready to fetch jobs", wn))

		j := <-q.Jobs
		ackd, _ := GetAck(q, hn, wn) // TODO: fix the current requeuing
		if !ackd {
			q.Jobs <- j
			return
		}
		b, err := json.Marshal(j)
		if err != nil {
			helper.FailOnError(err, "Could not decode job")
		}
		fmt.Fprintf(w, "%s", string(b))

		helper.ColorLog("\033[35m", fmt.Sprintf("Job:{Name:%s Args:%v} fetched by %s", j.JobName, j.Args, wn))

	}
}
