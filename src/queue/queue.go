package queue

import (
	"fmt"
	"log"
	"time"

	uuid "github.com/nu7hatch/gouuid"
)

//Queue ...
type Queue struct {
	Requests map[string]Request
}

//Request ...
type Request struct {
	ID     string
	Action func()
}

//NewQueue ...
func NewQueue() *Queue {
	q := &Queue{}
	q.Requests = make(map[string]Request)
	return q
}

//Run ...
func (q *Queue) Run() {

	for {
		if len(q.Requests) > 0 {
			log.Println("Performing...")
			q.PerformRequests()
		}
		time.Sleep(time.Second)
		log.Println("Queuing...")
	}

}

//PushRequest into queue
func (q *Queue) PushRequest(req Request) error {
	fmt.Println("Pushing requests...")
	u, err := uuid.NewV4()
	req.ID = u.String()
	q.Requests[req.ID] = req
	return err
}

//PerformRequests ...
func (q *Queue) PerformRequests() {
	rate := time.Second / 100
	burstLimit := 100
	tick := time.NewTicker(rate)
	defer tick.Stop()
	throttle := make(chan time.Time, burstLimit)
	go func() {
		for t := range tick.C {
			select {
			case throttle <- t:
			default:
			}
		} // exits after tick.Stop()
	}()
	i := len(q.Requests)
	for key, req := range q.Requests {
		<-throttle
		log.Printf("Go request[%d] %s\n", i, req.ID)
		go req.Action()
		delete(q.Requests, key)
		i = i - 1
	}
}
