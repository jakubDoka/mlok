package gate

import "github.com/jakubDoka/sterr"

var (
	errRunning    = sterr.New("calling g.%s while threads are still running, use g.Wait of g.CleanWait to synchronize")
	errNotRunning = sterr.New("calling g.%s while threads are not running, use g.Run before calling this")
)

// Gate supports concurent faze execution, you can break your game down to fazes that can be
// executed concurrently and put them into multiple Fazers
//
// for example if you need to process slice of entities and update of one entity does not mutate
// another, you can split processing between Fazers where each fazer will process n-th entity
type Gate struct {
	threads []*wrapper
	running bool
}

// alloc allocates new thread, old is reused if possible
func (g *Gate) alloc() (wp *wrapper) {
	l := len(g.threads)
	if cap(g.threads) != l {
		g.threads = g.threads[:l+1]
		if g.threads[l] == nil { // append can add more then cap + 1
			g.threads[l] = nWrapper()
		}
		wp = g.threads[l]
	} else {
		wp = nWrapper()
		g.threads = append(g.threads, wp)
	}
	wp.idx = l
	return
}

// Add adds Fazer to Gate and creates new goroutine that runs fazes
func (g *Gate) Add(fazer Fazer) {
	if g.running {
		panic(errRunning.Args("Add"))
	}

	wp := g.alloc()
	wp.Fazer = fazer

	go func() {
		for <-wp.input {
			f := wp.Faze()
			if f == nil {
				wp.output <- true
				wp.Fazer = nil
				return
			}
			f(wp.idx, len(g.threads))
			wp.output <- false
		}
		wp.output <- false
	}()
}

// Run runs all fazes, for each run, wait has to be called
func (g *Gate) Run() {
	if g.running {
		panic(errRunning.Args("Run"))
	}
	g.running = true

	for _, t := range g.threads {
		t.input <- true
	}
}

// Wait waits for threads to finish tasks
func (g *Gate) Wait() {
	if !g.running {
		panic(errNotRunning.Args("Wait"))
	}
	g.running = false

	for _, t := range g.threads {
		<-t.output
	}
}

// CleanWait cleans dead rotines, dead routine happens if Faze return nil
func (g *Gate) CleanWait() {
	if !g.running {
		panic(errNotRunning.Args("CleanWait"))
	}
	g.running = false

	var j, l int
	for i, k := 0, len(g.threads); i < k-j; i++ {
		t := g.threads[i]
		if <-t.output {
			j++
			if j == 1 {
				l = i
			}
			g.threads = append(append(g.threads[:i], g.threads[i+1:]...), t)[:k-j]
			t.Fazer = nil
		}
	}

	if j != 0 {
		for i := l; i < len(g.threads); i++ {
			g.threads[i].idx = i
		}
	}
}

// Clear disposes all threads, has to be called before dumping th Gate or memory leak will happen
func (g *Gate) Clear() {
	if g.running {
		panic(errRunning.Args("Clear"))
	}

	for _, t := range g.threads {
		t.input <- false
	}

	for _, t := range g.threads {
		<-t.output
	}

	g.threads = g.threads[:0]
}

// wrapper keeps additional info about fazer
type wrapper struct {
	input, output chan bool
	idx           int
	Fazer
}

func nWrapper() *wrapper {
	return &wrapper{
		input:  make(chan bool),
		output: make(chan bool),
	}
}

// FazeRunner is callback returned by Fazer.Faze that is then called by gate in async
type FazeRunner func(tIdx, count int)

// Fazer is something that can provide faze runner, when you register
// multiple threads to the gate at the same time gate will alway call one faze
// on all threads concurrently
type Fazer interface {
	// Faze can return nil, in that case Fazer will be disposed
	Faze() FazeRunner
}

// FazerBase is a most basic Fazer implementation, you can use ist as a base of your
// ovn thread types
type FazerBase struct {
	Fazes  []FazeRunner
	Cursor int
}

// Faze implements Fazer interface
func (t *FazerBase) Faze() FazeRunner {
	if len(t.Fazes) == 0 {
		return nil
	}
	i := t.Cursor
	t.Cursor = (t.Cursor + 1) % len(t.Fazes)
	return t.Fazes[i]
}
