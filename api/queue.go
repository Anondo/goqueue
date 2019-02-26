package api

import (
	"encoding/json"
	"fmt"
	"goqueue/helper"
	"goqueue/resources"
	"net/http"

	"github.com/go-chi/chi"
)

type QueueRequest struct {
	Name     string `json:"name"`
	Capacity int    `json:"cap"`
	Durable  bool   `json:"durable"`
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
		}

		resources.AddQueue(nq)
	}
}

type QueueListResponse struct {
	QNames []string `json:"qnames"`
}

func GetQueueList(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		qlr := QueueListResponse{}

		for _, q := range resources.QList {
			qlr.QNames = append(qlr.QNames, q.Name)
		}

		b, err := json.Marshal(qlr)

		if err != nil {
			helper.FailOnError(err, "Could not decode response")
		}

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
		} else {
			qdr.RMsg = "No Queue called: " + qn + " found"
		}

		b, err := json.Marshal(qdr)

		if err != nil {
			helper.FailOnError(err, "Could not decode response")
		}

		helper.ColorLog("\033[35m", fmt.Sprintf("Queue Deleted: %s\n", qn))

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
