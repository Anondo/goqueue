package api

import (
	"goqueue/helper"
	"goqueue/resources"
	"net/http"
	"strconv"
	"strings"
)

// SubscribeReq is the struct representing the request body for a subscribe request
type SubscribeReq struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	QName string `json:"qname"`
}

// SubscribeRequest is the handler api for subscribing a consumer against a queue
func SubscribeRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		sr := SubscribeReq{}

		helper.LogOnError(helper.ParseBody(r.Body, &sr), "Could not parse subscribe request")

		h := strings.Split(r.RemoteAddr, ":")[0]
		if h == "[" {
			h = "localhost"
		}

		p := strings.Split(r.RemoteAddr, ":")[1]
		pint, _ := strconv.Atoi(p)

		resources.SubscribeConsumer(h, pint, sr.QName, sr.Name, sr.ID)
	}
}

// UnsubscribeRequest is the http handler for unsubscribing a subscriber from the queue list
func UnsubscribeRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {

		ur := SubscribeReq{} // using the SubscribeReq here as well because both subscribe and unsubscribe have same payload

		helper.LogOnError(helper.ParseBody(r.Body, &ur), "Could not pare unsubscribe request")

		resources.UnsubscribeConsumer(ur.QName, ur.ID, ur.Name)
	}
}
