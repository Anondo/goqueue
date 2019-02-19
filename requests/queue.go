package requests

import (
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
