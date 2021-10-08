package spider

import (
	"net/http"
)

type Scheduler interface {
	Push(*http.Request)
	Poll() *http.Request

	Rest() int
}

type SchedulerChan struct {
	reqs chan *http.Request
}

func NewSchedulerChan() *SchedulerChan {
	reqs := make(chan *http.Request, 102400)
	return &SchedulerChan{reqs}
}

func (sc *SchedulerChan) Push(req *http.Request) {
	sc.reqs <- req
}

func (sc *SchedulerChan) Poll() *http.Request {
	if len(sc.reqs) == 0 {
		return nil
	}
	return <-sc.reqs
}

func (sc *SchedulerChan) Rest() int {
	return len(sc.reqs)
}
