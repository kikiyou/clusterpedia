package queue

import (
	"sync"
	"time"

	"k8s.io/utils/clock"
)

// WorkQueue allows queuing items with a timestamp. An item is
// considered ready to process if the timestamp has expired.
type WorkQueue interface {
	// GetWork dequeues and returns all ready items.
	GetWork() []string
	// Enqueue inserts a new item or overwrites an existing item.
	Enqueue(item string, delay time.Duration)
	Len() int
}

type basicWorkQueue struct {
	clock clock.Clock
	lock  sync.Mutex
	queue map[string]time.Time
}

var _ WorkQueue = &basicWorkQueue{}

func NewBasicWorkQueue(clock clock.Clock) WorkQueue {
	queue := make(map[string]time.Time)
	return &basicWorkQueue{queue: queue, clock: clock}
}

func (q *basicWorkQueue) GetWork() []string {
	q.lock.Lock()
	defer q.lock.Unlock()
	now := q.clock.Now()
	var items []string
	for k, v := range q.queue {
		if v.Before(now) {
			items = append(items, k)
			delete(q.queue, k)
		}
	}
	return items
}

func (q *basicWorkQueue) Enqueue(item string, delay time.Duration) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue[item] = q.clock.Now().Add(delay)
}

func (q *basicWorkQueue) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.Len()
}
