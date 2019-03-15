package api

import (
	"goqueue/helper"
	"goqueue/resources"
	"net/http"
	"strconv"
	"strings"
)

type SubscribeResp struct {
	Name  string `json:"name"`
	QName string `json:"qname"`
}

func SubscribeRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		sr := SubscribeResp{}

		helper.LogOnError(helper.ParseBody(r.Body, &sr), "Could not parse subscribe request")

		h := strings.Split(r.RemoteAddr, ":")[0]
		if h == "[" {
			h = "localhost"
		}

		p := strings.Split(r.RemoteAddr, ":")[1]
		pint, _ := strconv.Atoi(p)

		resources.SubscribeConsumer(h, pint, sr.QName, sr.Name)
	}
}
