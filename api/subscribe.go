package api

import (
	"goqueue/helper"
	"goqueue/resources"
	"net/http"
)

type SubscribeResp struct {
	Name  string `json:"name"`
	Port  int    `json:"port"`
	QName string `json:"qname"`
}

func SubscribeRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		sr := SubscribeResp{}

		helper.LogOnError(helper.ParseBody(r.Body, &sr), "Could not parse subscribe request")

		h := r.URL.Hostname()
		if h == "" {
			h = "localhost"
		}

		resources.SubscribeConsumer(h, sr.Port, sr.QName, sr.Name)
	}
}
