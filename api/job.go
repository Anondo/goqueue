package api

import (
	"goqueue/helper"
	"goqueue/resources"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spf13/viper"
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
			j.QName = viper.GetString("default.queue_name")
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

	resources.SendJob(w, qn, wn)

}
