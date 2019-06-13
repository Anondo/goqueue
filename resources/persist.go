package resources

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

var (
	durableFileName string
)

// JSONJob is the correspondent struct of the JOB struct
type JSONJob struct {
	ID      int         `json:"id"`
	JobName string      `json:"job_name"`
	Args    []Arguments `json:"args"`
}

// JSONQueue is the correspondent struct of the Queue struct
type JSONQueue struct {
	ID                  int           `json:"id"`
	Name                string        `json:"name"`
	Jobs                []Job         `json:"jobs"`
	Capacity            int           `json:"capacity"`
	RegisteredTaskNames []string      `json:"reg_tasks"`
	Subscribers         []*Subscriber `json:"subscribers"`
	AckWait             time.Duration `json:"ack"`
	Durable             bool          `json:"durable"`
}

// FromJSON Populates a queue from a JSONQueue
func (jq *JSONQueue) FromJSON() Queue {
	jc := make(chan Job, jq.Capacity)

	for _, j := range jq.Jobs {
		jc <- Job{
			ID:      j.ID,
			JobName: j.JobName,
			Args:    j.Args,
		}
	}

	return Queue{
		ID:                  jq.ID,
		Name:                jq.Name,
		Jobs:                jc,
		Capacity:            jq.Capacity,
		RegisteredTaskNames: jq.RegisteredTaskNames,
		Subscribers:         jq.Subscribers,
		AckWait:             jq.AckWait,
		Durable:             jq.Durable,
	}
}

// ToJSON populates a JSONQueue from Queue
func (q *Queue) ToJSON() JSONQueue {
	return JSONQueue{
		ID:                  q.ID,
		Name:                q.Name,
		Jobs:                []Job{},
		Capacity:            q.Capacity,
		RegisteredTaskNames: q.RegisteredTaskNames,
		Subscribers:         q.Subscribers,
		AckWait:             q.AckWait,
		Durable:             q.Durable,
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

	jq := q.ToJSON()

	jql = append(jql, jq)

	b, _ := json.Marshal(jql)
	return ioutil.WriteFile(durableFileName, b, 0644)

}

func (q *Queue) addDurableSubscriber(s *Subscriber) error {
	data, err := ioutil.ReadFile(durableFileName)

	jql := []*JSONQueue{}

	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &jql); err != nil {
		return err
	}

	for _, jq := range jql {
		if jq.Name == q.Name {
			jq.Subscribers = append(jq.Subscribers, s)
			b, _ := json.Marshal(jql)
			return ioutil.WriteFile(durableFileName, b, 0644)
		}
	}
	return nil
}

func (q *Queue) removeDurableSubscriber(s *Subscriber) error {
	data, err := ioutil.ReadFile(durableFileName)

	jql := []*JSONQueue{}

	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &jql); err != nil {
		return err
	}

	for _, jq := range jql {
		if jq.Name == q.Name {
			for i, sbscrbr := range jq.Subscribers {
				if sbscrbr.ID == s.ID {
					jq.Subscribers = append(jq.Subscribers[:i], jq.Subscribers[i+1:]...)
					b, _ := json.Marshal(jql)
					return ioutil.WriteFile(durableFileName, b, 0644)
				}
			}
		}
	}
	return nil
}

func (q *Queue) addDurableRegTask(tsk string) error {
	data, err := ioutil.ReadFile(durableFileName)

	jql := []*JSONQueue{}

	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &jql); err != nil {
		return err
	}

	for _, jq := range jql {
		if jq.Name == q.Name {
			jq.RegisteredTaskNames = append(jq.RegisteredTaskNames, tsk)
			b, _ := json.Marshal(jql)
			return ioutil.WriteFile(durableFileName, b, 0644)
		}
	}
	return nil
}

func (q *Queue) addDurableJob(j Job) error {
	data, err := ioutil.ReadFile(durableFileName)

	jql := []*JSONQueue{}

	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &jql); err != nil {
		return err
	}

	for _, jq := range jql {
		if jq.Name == q.Name {
			jq.Jobs = append(jq.Jobs, j)
			b, _ := json.Marshal(jql)
			return ioutil.WriteFile(durableFileName, b, 0644)
		}
	}
	return nil
}

func (q *Queue) removeDurableJob(j Job) error {
	data, err := ioutil.ReadFile(durableFileName)

	jql := []*JSONQueue{}

	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &jql); err != nil {
		return err
	}

	for _, jq := range jql {
		if jq.Name == q.Name {
			jq.Jobs = jq.Jobs[:len(jq.Jobs)-1]
			b, _ := json.Marshal(jql)
			return ioutil.WriteFile(durableFileName, b, 0644)
		}
	}
	return nil
}

// RemovePersistedQueue removes a queue from persistance
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
			b, _ := json.Marshal(jql)
			return ioutil.WriteFile(durableFileName, b, 0644)
		}
	}
	return nil
}

// Requeue re-pushes a job to the queue, used in case of some kind of failure
func (q *Queue) Requeue() error {
	data, err := ioutil.ReadFile(durableFileName)

	jql := []JSONQueue{}

	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &jql); err != nil {
		return err
	}

	for _, jq := range jql {
		if q.Name == jq.Name {
			q.Jobs = jq.FromJSON().Jobs
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
