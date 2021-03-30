package frame

import (
	"fmt"
	"time"
)

// Delta is delta time messuring
type Delta struct {
	time.Time
	fps  int
	time float64
}

// Init initializes delta to get rid of firs very long frame
func (d Delta) Init() Delta {
	d.Time = time.Now()
	return d
}

// Tick updates delta
func (d *Delta) Tick() float64 {
	delta := time.Since(d.Time).Seconds()
	d.time += delta
	d.fps++
	d.Time = time.Now()
	return delta
}

// Log logs fps to console every interval (in seconds)
func (d *Delta) Log(interval float64) {
	d.CustomLog(interval, nil)
}

// Log logs fps to console every interval (in seconds), you can optionally log more
// information with runner
func (d *Delta) CustomLog(interval float64, runner func()) {
	if d.time > interval {
		fmt.Println("fps:", float64(d.fps)/interval)
		d.time = 0
		d.fps = 0
		if runner != nil {
			runner()
		}
	}
}

// Limitter controls frame count per second, by deafult limmiter does nothing
// fps has to be set
type Limitter struct {
	now   time.Time
	frame time.Duration
}

// SetFPS sets fps to be limited to a given value
func (l *Limitter) SetFPS(fps int) {
	l.frame = time.Second / time.Duration(fps)
	l.now = time.Now()
}

// Regulate performs frame regulation (call every frame)
func (l *Limitter) Regulate() {
	sleepTime := l.frame - time.Duration(time.Since(l.now).Nanoseconds())
	time.Sleep(sleepTime)
	l.now = l.now.Add(l.frame)
}
