package resources

import (
	"goqueue/helper"
)

type Subscriber struct {
	Host  string `json:"host"`
	Port  int    `json:"port"`
	CName string `json:"cname"`
	Ack   bool   `json:"ack"`
}

func SubscribeConsumer(h string, p int, qn, wn string) {
	q := GetQueueByName(qn)

	if q != nil {
		for _, s := range q.Subscribers {
			if s.Host == h && s.Port == p && s.CName == wn {
				return
			}
		}
		s := &Subscriber{
			Host:  h,
			Port:  p,
			CName: wn,
		}
		q.Subscribers = append(q.Subscribers, s)
		helper.FailOnError(q.addDurableSubscriber(s), "Failed to persist subscriber")
	}

	helper.ColorLog("\033[35m", "Successfully subscribed:"+wn)

}
