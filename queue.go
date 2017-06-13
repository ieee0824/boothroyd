package boothroyd

import (
	"time"
	"github.com/ieee0824/getenv"
)

var (
	MAX_QUEUE_SIZE = getenv.Int("MAX_QUEUE_SIZE", 1 * 1024 * 1024)
)

type innerQueue struct {
	data chan interface{}
	status []interface{}
	c chan interface{}
	delay time.Duration
	lastDeque time.Time
}

func newInnnerQueue() *innerQueue {
	return &innerQueue {
		make(chan interface{}, MAX_QUEUE_SIZE),
		[]interface{}{},
		make(chan interface{}),
		time.Duration(int64(getenv.Int("DELAY_TIME", 5))) * time.Second,
		time.Now().AddDate(-1,0,0),
	}
}

func (q *innerQueue) enqueue(i interface{}) {
	q.data <- i
	q.status = append(q.status, i)
}

func (q *innerQueue) dequeue() chan interface{} {
	go func() {
		for {
			d := <- q.data

			for time.Now().Sub(q.lastDeque) < q.delay {
				time.Sleep(q.delay - time.Now().Sub(q.lastDeque))
			}
			q.c <- d
			q.status = q.status[1:]
			q.lastDeque = time.Now()
		}
	}()
	return q.c
}

type Queue struct {
	C chan interface{}
	data map[string]*innerQueue
}

var Desmond = New

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

func (q *Queue) Dequeue() interface{} {
	return <- q.C
}

func (q *Queue) Status() map[string]interface{} {
	var r = map[string]interface{}{}
	for k, v := range q.data {
		r[k] = map[string]interface{}{
			"len": len(v.status),
			"targets": v.status,
		}
	}
	return r
}