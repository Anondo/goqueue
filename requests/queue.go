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
	if r.Method == "PUT" {
		qr := QueueRequest{}
		helper.FailOnError(helper.ParseBody(r.Body, &qr), "Could not parse Queue Info")

		resources.Q.Name = qr.Name
		resources.Q.Capacity = qr.Capacity

		log.Printf("Queue Declared With Name:%s & Capacity:%d\n", resources.Q.Name, resources.Q.Capacity)
	}
}
