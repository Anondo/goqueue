package requests

import (
	"goqueue/helper"
	"goqueue/resources"
	"net/http"
)

type JobCreateRequest struct {
	Task string        `json:"task"`
	Args []interface{} `json:"args"`
}

func CreateJobRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		j := JobCreateRequest{}
		helper.FailOnError(helper.ParseBody(r.Body, &j), "Failed to Create Job")

		resources.Q.PushTask(j.Task, j.Args)
	}
}
