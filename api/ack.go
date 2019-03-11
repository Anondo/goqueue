package api

import (
	"fmt"
	"goqueue/helper"
	"goqueue/resources"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

type AckRequest struct {
	Ack        bool   `json:"ack"`
	Qname      string `json:"qname"`
	Subscriber string `json:"subscriber"`
}

func Acknowledgement(w http.ResponseWriter, r *http.Request) {
	a := AckRequest{}
	helper.LogOnError(helper.ParseBody(r.Body, &a), "Couldn't parse ackonwledgement request body")

	q := resources.GetQueueByName(a.Qname)

	fmt.Println("heeeeeerrreeeeeeee")

	if q != nil {
		s := q.GetSubscriber(a.Subscriber)
		if s != nil {
			s.Ack = a.Ack
			spew.Dump(s)
		}
	}

}
