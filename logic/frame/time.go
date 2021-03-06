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
	if d.time > interval {
		fmt.Println("fps:", float64(d.fps)/interval)
		d.time = 0
		d.fps = 0
	}
}
