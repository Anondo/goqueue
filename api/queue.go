package api

import (
	"encoding/json"
	"fmt"
	"goqueue/helper"
	"goqueue/resources"
	"net/http"
)

type QueueRequest struct {
	Name     string `json:"name"`
	Capacity int    `json:"cap"`
}

func DeclearQueue(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
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
		}

		resources.AddQueue(nq)

		helper.ColorLog("\033[35m", fmt.Sprintf("Queue Declared: {Name:%s & Capacity:%d}\n", qr.Name, qr.Capacity))
	}
}

type QueueListResponse struct {
	QNames []string `json:"qnames"`
}

func GetQueueList(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
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
