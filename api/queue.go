package api

import (
	"encoding/json"
	"fmt"
	"goqueue/helper"
	"goqueue/resources"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type QueueRequest struct {
	Name     string        `json:"name"`
	Capacity int           `json:"cap"`
	Durable  bool          `json:"durable"`
	AckWait  time.Duration `json:"ack_wait"`
}

func DeclearQueue(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		qr := QueueRequest{}
		helper.LogOnError(helper.ParseBody(r.Body, &qr), "Could not parse Queue Info")

		if qr.Name == "" {
			return
		}

		nq := resources.Queue{
			ID:       len(resources.QList) + 1,
			Name:     qr.Name,
			Capacity: qr.Capacity,
			Jobs:     make(chan resources.Job, qr.Capacity),
			Durable:  qr.Durable,
			AckWait:  qr.AckWait,
		}

		resources.AddQueue(nq)
	}
}

type QueueListResponse struct {
	Queues []resources.JSONQueue `json:"queues"`
}

func GetQueueList(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		qlr := QueueListResponse{}

		for _, q := range resources.QList {
			qlr.Queues = append(qlr.Queues, q.ToJSON())
		}

		b, err := json.Marshal(qlr)

		if err != nil {
			helper.FailOnError(err, "Could not decode response")
		}

		fmt.Fprintf(w, "%s", string(b))
	}
}

func ClearQueue(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		qn := chi.URLParam(r, "qname")
		q := resources.GetQueueByName(qn)

		var qcr struct {
			RMsg string `json:"response_message"`
		}

		if q == nil {
			qcr.RMsg = "No queue named: " + qn + " was found"
		} else {
			qcr.RMsg = q.Clear()
		}

		b, _ := json.Marshal(qcr)

		fmt.Fprintf(w, "%s", string(b))
	}
}

type QueueDeleteResponse struct {
	RMsg string `json:"response_message"`
}

func DeleteQueue(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		qn := chi.URLParam(r, "qname")

		found := false
		for i, q := range resources.QList {
			if q.Name == qn {
				resources.QList = append(resources.QList[:i], resources.QList[i+1:]...)
				resources.RemovePersistedQueue(qn)
				found = true
				break
			}
		}
		qdr := QueueDeleteResponse{}
		if found {
			qdr.RMsg = "Queue Deleted Successfully"
			helper.ColorLog("\033[35m", fmt.Sprintf("Queue Deleted: %s\n", qn))
		} else {
			qdr.RMsg = "No Queue called: " + qn + " found"
		}

		b, err := json.Marshal(qdr)

		if err != nil {
			helper.FailOnError(err, "Could not decode response")
		}

		fmt.Fprintf(w, "%s", string(b))

	}
}

type RegTaskReq struct {
	TaskNames []string `json:"task_names"`
	QName     string   `json:"qname"`
}

func RegisterTaskRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		rtr := RegTaskReq{}
		if err := helper.ParseBody(r.Body, &rtr); err != nil {
			helper.LogOnError(err, "Could not parse registration request body")
		}

		if err := resources.RegisterTasks(rtr.QName, rtr.TaskNames); err != nil {
			helper.LogOnError(err, "Failed to register tasks")
		}

	}
}
