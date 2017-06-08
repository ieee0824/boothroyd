package queue

import (
	"time"
)

const (
	MAX_QUEUE_SIZE = 1 *1024 * 1024
)

type innerQueue struct {
	data chan interface{}
	status int
	c chan interface{}
	delay time.Duration
	lastDeque time.Time
}

func newInnnerQueue() *innerQueue {
	return &innerQueue {
		make(chan interface{}, MAX_QUEUE_SIZE),
		0,
		make(chan interface{}),
		5 * time.Second,
		time.Now().AddDate(-1,0,0),
	}
}

func (q *innerQueue) enqueue(i interface{}) {
	q.data <- i
	q.status ++
}

func (q *innerQueue) dequeue() chan interface{} {
	go func() {
		for {
			d := <- q.data

			if time.Now().Sub(q.lastDeque) < q.delay {
				for time.Now().Sub(q.lastDeque) < q.delay {
					time.Sleep(q.delay - time.Now().Sub(q.lastDeque))
				}
			}
			q.c <- d
			q.status --
			q.lastDeque = time.Now()
		}
	}()
	return q.c
}

type Queue struct {
	C chan interface{}
	data map[string]*innerQueue
}

func New()*Queue {
	return &Queue{
		make(chan interface{}),
		map[string]*innerQueue{},
	}
}

func (q *Queue) Enqueue(key string, i interface{}) {
	if _, ok := q.data[key]; !ok {
		q.data[key] = newInnnerQueue()

	}
	q.data[key].enqueue(i)
	go func() {
		for {
			select {
			case d := <- q.data[key].dequeue():
				q.C <- d
			}
		}
	}()
}

func (q *Queue) Status() map[string]interface{} {
	var r = map[string]interface{}{}
	for k, v := range q.data {
		r[k] = v.status
		r[k] = v.status
	}
	return r
}