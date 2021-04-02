package timer

import "math/rand"

// Timer measures a elapsed time, its very simple but effective
type Timer struct {
	Progress, Period float64
}

// Period returns timer with given period
func Period(period float64) Timer {
	return Timer{Period: period}
}

// Progress returns timer with given progress and period
func Progress(progress float64, period float64) Timer {
	return Timer{progress, period}
}

func Random(period float64) Timer {
	return Timer{period * rand.Float64(), period}
}

// Tick increases progress by delta
func (t *Timer) Tick(delta float64) {
	t.Progress += delta
}

// Done returns whether progress exceeded period
func (t Timer) Done() bool {
	return t.Progress >= t.Period
}

// Reset sets progres to 0
func (t *Timer) Reset() {
	t.Progress = 0
}

// Skip skips timer cycle
func (t *Timer) Skip() {
	t.Progress = t.Period
}

// DoneReset returns true and resets the timer
func (t *Timer) DoneReset() bool {
	if t.Progress >= t.Period {
		t.Progress = 0
		return true
	}
	return false
}

// TickDone does Tick and Done in one step
func (t *Timer) TickDone(delta float64) bool {
	t.Progress += delta
	return t.Progress >= t.Period
}

// TickDoneReset does tick and if Done that it resets
func (t *Timer) TickDoneReset(delta float64) bool {
	t.Progress += delta
	if t.Progress > t.Period {
		t.Progress = 0
		return true
	}
	return false
}
