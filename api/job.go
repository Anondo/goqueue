package api

import (
	"goqueue/helper"
	"goqueue/resources"
	"net/http"

	"github.com/go-chi/chi"
)

type JobCreateRequest struct {
	Task  string                `json:"task"`
	Args  []resources.Arguments `json:"args"`
	QName string                `json:"qname"`
}

func CreateJobRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
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

func FetchJobRequest(w http.ResponseWriter, r *http.Request) {
	qn := chi.URLParam(r, "queue_name")
	wn := r.URL.Query().Get("sname")

	if qn == "" || !resources.QueueExists(qn) {
		return
	}

	hn := r.URL.Hostname()
	if hn == "" {
		hn = "localhost"
	}

	resources.SendJob(w, qn, wn, hn)

}
