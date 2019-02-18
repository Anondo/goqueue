package requests

import (
	"goqueue/helper"
	"goqueue/resources"
	"log"
	"net/http"
)

type QueueRequest struct {
	Name     string `json:"name"`
	Capacity int    `json:"cap"`
}

func DeclearQueue(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		qr := QueueRequest{}
		helper.FailOnError(helper.ParseBody(r.Body, &qr), "Could not parse Queue Info")

		nq := resources.Queue{
			ID:       len(resources.QList) + 1,
			Name:     qr.Name,
			Capacity: qr.Capacity,
			Jobs:     make(chan resources.Job, qr.Capacity),
		}

		resources.AddQueue(nq)

		log.Printf("Queue Declared With Name:%s & Capacity:%d\n", qr.Name, qr.Capacity)
	}
}
