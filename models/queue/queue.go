package queue

import (
	"time"
	"github.com/pkg/errors"
)

type Queue struct {
	LastPopedAt *time.Time
	Weight time.Duration
	Data []string
}

func New()*Queue{
	return &Queue{
		nil,
		10 * time.Second,
		make([]string, 0, 1024),
	}
}

func (q *Queue) Enqueue (dataSet ...string) {
	q.Data = append(q.Data, dataSet...)
}

func (q *Queue) Dequeue () (string, error) {
	if q == nil {
		return "", errors.New("queue is nil")
	}
	if len(q.Data) == 0 {
		return "", errors.New("queue is empty")
	}

	if q.LastPopedAt != nil && (time.Now().Sub(*q.LastPopedAt) < q.Weight) {
		return "", errors.New("Dequeue Time interval too short")
	}

	ret := q.Data[0]
	q.Data = q.Data[1:]
	now := time.Now()
	q.LastPopedAt = &now
	return ret, nil
}

