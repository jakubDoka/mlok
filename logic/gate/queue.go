package gate

import "sync"

// Queue can be used to push tasks between threads.
// mind that even though this can simplify things greatly
// pushing func() with captured stack will produce allocations
// and mutex also can slow things down
type Queue struct {
	m sync.Mutex

	tasks, help []func()
}

// Post posts task to queue, can be called from anywhere
func (q *Queue) Post(task func()) {
	q.m.Lock()
	q.tasks = append(q.tasks, task)
	q.m.Unlock()
}

// Run runs all tasks and clears the queue
// should be called from thread you want to execute
// the tasks in
func (q *Queue) Run() {
	q.m.Lock()
	q.tasks, q.help = q.help, q.tasks
	q.m.Unlock()

	for _, t := range q.help {
		t()
	}
	q.help = q.help[:0]
}
