package requests

import (
	"goqueue/helper"
	"goqueue/resources"
	"net/http"
)

type JobCreateRequest struct {
	Task  string        `json:"task"`
	Args  []interface{} `json:"args"`
	QName string        `json:"qname"`
}

func CreateJobRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		j := JobCreateRequest{}
		helper.LogOnError(helper.ParseBody(r.Body, &j), "Failed to Create Job")

		if j.QName == "" {
			j.QName = "default_queue"
		}

		if j.Task == "" {
			return
		}

		resources.AddTask(j.QName, j.Task, j.Args)
	}
}
