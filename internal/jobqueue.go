package internal

import "sync"

// JobQueue is a FIFO queue of Jobs. Each repository has its own queue to ensure that jobs
// are run in order one at a time.
type JobQueue struct {
	items []*Job
	lock  *sync.Mutex
}

func NewJobQueue() *JobQueue {
	return &JobQueue{items: make([]*Job, 0), lock: &sync.Mutex{}}
}

func (q *JobQueue) Add(job *Job) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.items = append(q.items, job)
}

func (q *JobQueue) Pop() *Job {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.items) > 0 {
		x := q.items[0]
		q.items = q.items[1:]
		return x
	}
	return nil
}

func (q *JobQueue) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return len(q.items)
}
