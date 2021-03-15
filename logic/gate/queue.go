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

// Worker is simple reusable thread for general use
// it can work only on one task and tasks cannot stack
// up, after each Do, Wait has to be called. Calling any worker
// method inside his task will result in deadlock inside worker thread
type Worker struct {
	working bool
	in      chan func()
	out     chan struct{}
}

// NWorker sets up a thread on witch the worker is working on
func NWorker() Worker {
	in, out := make(chan func()), make(chan struct{})

	go func() {
		for task := <-in; task != nil; task = <-in {
			task()
			out <- struct{}{}
		}
	}()

	return Worker{false, in, out}
}

// Kill makes worker thread dispatch
//
// panics in same cases as Do
func (w Worker) Kill() {
	w.Do(nil) // poetical
}

// Do sends a task to worker to work on, if you pass nil, thread is dispatched
//
// panics if worker is already working, to prevent this call Wait
func (w *Worker) Do(task func()) {
	if w.working {
		panic("worker is already working on someting, call wait to synchronize")
	}
	w.working = true
	w.in <- task
}

// Wait waits for worker to finish a job
//
// panics if worker is not working, this happens if you call Wait twice without Do in between
func (w *Worker) Wait() {
	if !w.working {
		panic("worker is not working, Wait would block a this thread")
	}
	<-w.out
	w.working = false
}
