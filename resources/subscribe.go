package resources

import "goqueue/helper"

type Subscriber struct {
	Host  string
	Port  int
	CName string
}

func SubscribeConsumer(h string, p int, qn, wn string) {
	q := GetQueueByName(qn)

	if q != nil {
		for _, s := range q.Subscribers {
			if s.Host == h && s.Port == p && s.CName == wn {
				return
			}
		}
		q.Subscribers = append(q.Subscribers, &Subscriber{
			Host:  h,
			Port:  p,
			CName: wn,
		})
	}

	helper.ColorLog("\033[35m", "Successfully subscribed:"+wn)

}
