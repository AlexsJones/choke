package queue

import (
	"fmt"
	"log"
	"sync"
	"time"

	uuid "github.com/nu7hatch/gouuid"
)

//Queue ...
type Queue struct {
	Requests map[string]Request
	lock     *sync.RWMutex
}

//Request ...
type Request struct {
	ID     string
	Action func()
}

//NewQueue ...
func NewQueue() *Queue {
	q := &Queue{}
	q.lock = &sync.RWMutex{}
	q.Requests = make(map[string]Request)
	return q
}

//Run ...
func (q *Queue) Run() {

	for {
		if len(q.Requests) > 0 {
			q.PerformRequests()
			log.Printf("Chunked queue remaining %d\n", len(q.Requests))
		}
		time.Sleep(time.Second)
	}

}

//PushRequest into queue
func (q *Queue) PushRequest(req Request) error {
	fmt.Println("Pushing requests...")
	u, err := uuid.NewV4()
	req.ID = u.String()
	q.lock.Lock()
	defer q.lock.Unlock()
	q.Requests[req.ID] = req
	return err
}

//PerformRequests ...
func (q *Queue) PerformRequests() {
	rate := time.Second / 10
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
		}
	}()
	i := len(q.Requests)
	for key, req := range q.Requests {
		<-throttle
		go req.Action()
		q.lock.Lock()
		defer q.lock.Unlock()
		delete(q.Requests, key)
		i = i - 1
	}
}
