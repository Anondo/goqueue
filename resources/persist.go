package resources

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var (
	durableFileName string
)

type JSONJob struct {
	ID      int         `json:"id"`
	JobName string      `json:"job_name"`
	Args    []Arguments `json:"args"`
}

type JSONQueue struct {
	ID                  int           `json:"id"`
	Name                string        `json:"name"`
	Jobs                []JSONJob     `json:"jobs"`
	Capacity            int           `json:"capacity"`
	RegisteredTaskNames []string      `json:"reg_tasks"`
	Subscribers         []*Subscriber `json:"subscribers"`
}

func (j *JSONQueue) FromJSON() Queue {
	jc := make(chan Job, j.Capacity)

	for _, j := range j.Jobs {
		jc <- Job{
			ID:      j.ID,
			JobName: j.JobName,
			Args:    j.Args,
		}
	}

	return Queue{
		ID:                  j.ID,
		Name:                j.Name,
		Jobs:                jc,
		Capacity:            j.Capacity,
		RegisteredTaskNames: j.RegisteredTaskNames,
		Subscribers:         j.Subscribers,
	}
}

func (q *Queue) toJSON() JSONQueue {
	return JSONQueue{
		ID:                  q.ID,
		Name:                q.Name,
		Jobs:                []JSONJob{},
		Capacity:            q.Capacity,
		RegisteredTaskNames: q.RegisteredTaskNames,
		Subscribers:         q.Subscribers,
	}
}

func (q *Queue) persistQueue() error {
	data, err := ioutil.ReadFile(durableFileName)

	jql := []JSONQueue{}

	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &jql); err != nil {
		return err
	}

	jq := q.toJSON()

	jql = append(jql, jq)

	json_file, err := os.OpenFile(durableFileName, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	if err := json.NewEncoder(json_file).Encode(jql); err != nil {
		return err
	}
	return json_file.Close()

}

func RemovePersistedQueue(qn string) error {
	data, err := ioutil.ReadFile(durableFileName)

	jql := []JSONQueue{}

	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &jql); err != nil {
		return err
	}

	for i, jq := range jql {
		if jq.Name == qn {
			jql = append(jql[:i], jql[i+1:]...)
			json_file, err := os.OpenFile(durableFileName, os.O_RDWR, 0644)
			if err != nil {
				return err
			}
			if err := json.NewEncoder(json_file).Encode(jql); err != nil {
				return err
			}
			return json_file.Close()
		}
	}
	return nil
}

func initPersistedQueues() error {
	data, err := ioutil.ReadFile(durableFileName)

	jql := []JSONQueue{}

	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &jql); err != nil {
		return err
	}

	qs := []Queue{}

	for _, jq := range jql {
		qs = append(qs, jq.FromJSON())
	}

	QList = append(QList, qs...)

	return nil
}
