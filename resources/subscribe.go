package resources

import (
	"goqueue/helper"
)

// Subscriber represents a subscriber informations
type Subscriber struct {
	ID    string
	Host  string `json:"host"`
	Port  int    `json:"port"`
	CName string `json:"cname"`
	Ack   bool   `json:"ack"`
}

// SubscribeConsumer subscribes a consumer for a particular queue
func SubscribeConsumer(h string, p int, qn, wn, id string) {
	q := GetQueueByName(qn)

	if q != nil {
		for _, s := range q.Subscribers {
			if s.ID == id {
				return
			}
		}
		s := &Subscriber{
			ID:    id,
			Host:  h,
			Port:  p,
			CName: wn,
		}
		q.Subscribers = append(q.Subscribers, s)
		helper.FailOnError(q.addDurableSubscriber(s), "Failed to persist subscriber")
	}

	helper.ColorLog("\033[35m", "Successfully subscribed:"+wn)

}

// func UnsubscribeConsumer(id string)
