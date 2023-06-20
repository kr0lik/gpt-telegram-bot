package queue

import (
	"sync"
	"time"
)

const (
	rateLimitTime  = 60
	rateLimitQuery = 3
)

type OpenAi struct {
	queue []time.Time
	mu    *sync.Mutex
}

func NewOpenAi() *OpenAi {
	return &OpenAi{mu: &sync.Mutex{}}
}

func (q *OpenAi) IsLocked() (bool, time.Duration) {
	q.mu.Lock()
	defer q.mu.Unlock()

	lock, expire := q.check()
	if lock {

		return lock, expire
	}

	q.queue = append(q.queue, time.Now())

	return lock, 0
}

func (q *OpenAi) check() (bool, time.Duration) {
	now := time.Now()
	count := 0
	expire := 0

	queue := make([]time.Time, 0)

	for _, item := range q.queue {
		itemTime := int(now.Sub(item).Seconds())

		if itemTime <= rateLimitTime {
			if expire == 0 || itemTime > expire {
				expire = itemTime
			}

			queue = append(queue, item)

			count++
		}
	}

	q.queue = queue

	return count >= rateLimitQuery, time.Duration(rateLimitTime-expire) * time.Second
}
