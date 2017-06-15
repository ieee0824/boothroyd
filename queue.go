package boothroyd

import (
	"time"
	"github.com/ieee0824/getenv"
	"encoding/json"
)

var (
	MAX_QUEUE_SIZE = getenv.Int("MAX_QUEUE_SIZE", 1 * 1024 * 1024)
	GC_PARAM = time.Duration(int64(getenv.Int("GC_PARAM", 10))) * time.Second
)

type innerQueue struct {
	data chan interface{}
	status []interface{}
	c chan interface{}
	delay time.Duration
	lastDeque time.Time
	lockFlag bool
}

func newInnnerQueue() *innerQueue {
	return &innerQueue {
		make(chan interface{}, MAX_QUEUE_SIZE),
		[]interface{}{},
		make(chan interface{}),
		time.Duration(int64(getenv.Int("DELAY_TIME", 5))) * time.Second,
		time.Now().AddDate(-1,0,0),
		false,
	}
}

func (q *innerQueue) lock() {
	q.lockFlag = true
}

func (q *innerQueue) unlock() {
	q.lockFlag = false
}

func (q *innerQueue) enqueue(i interface{}) {
	q.status = append(q.status, i)
	q.data <- i
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
	q :=  &Queue{
		make(chan interface{}),
		map[string]*innerQueue{},
	}

	q.gc()

	return q
}

func (q *Queue) gc() {
	t := time.NewTicker(GC_PARAM)

	go func() {
		for {
			select {
			case <-t.C:
				for k, _ := range q.data {
					if len(q.data[k].status) == 0 {
						q.data[k].lock()
					} else {
						continue
					}
					if GC_PARAM * 2 < time.Now().Sub(q.data[k].lastDeque) {
						delete(q.data, k)
					} else {
						q.data[k].unlock()
					}
				}
			}
		}
	}()
}

func (q *Queue) Enqueue(key string, i interface{}) {
	retry:
	if _, ok := q.data[key]; !ok {
		q.data[key] = newInnnerQueue()
	}
	if _, ok := q.data[key]; !ok || q.data[key].lockFlag {
		goto retry
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

func (q *Queue) IsEmpty() bool {
	if len(q.data) == 0 {
		return true
	}
	for _, v := range q.data {
		if len(v.status) != 0 {
			return false
		}
	}
	return true
}

func (q Queue) String() string {
	bin, err := json.Marshal(q.Status())
	if err != nil {
		return err.Error()
	}
	return string(bin)
}