package api

import (
	"goqueue/helper"
	"goqueue/resources"
	"net/http"
)

// AckRequest is the struct representing the request payload for Acknowledgement
type AckRequest struct {
	Ack        bool   `json:"ack"`
	Qname      string `json:"qname"`
	Subscriber string `json:"subscriber"`
}

// Acknowledgement is the handler api for getting ackonwledgement from the consumer
func Acknowledgement(w http.ResponseWriter, r *http.Request) {
	a := AckRequest{}
	helper.LogOnError(helper.ParseBody(r.Body, &a), "Couldn't parse ackonwledgement request body")

	q := resources.GetQueueByName(a.Qname)

	if q != nil {
		s := q.GetSubscriber(a.Subscriber)
		if s != nil {
			s.Ack = a.Ack
		}
	}

}
