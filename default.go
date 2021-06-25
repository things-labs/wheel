package wheel

import (
	"sync"
	"time"
)

var once sync.Once
var base = New()

func lazyInit() { once.Do(func() { base.Run() }) }

// HasRunning base running status.
func HasRunning() bool {
	return base.HasRunning()
}

// Len the number timer of the base.
func Len() int {
	lazyInit()
	return base.Len()
}

// AddJob add a job
func AddJob(job Job, timeout time.Duration) *Timer {
	lazyInit()
	return base.AddJob(job, timeout)
}

// AddJobFunc add a job function
func AddJobFunc(f func(), timeout time.Duration) *Timer { return AddJob(JobFunc(f), timeout) }

// Add add timer to base. and start immediately.
func Add(tm *Timer, timeout time.Duration) {
	lazyInit()
	base.Add(tm, timeout)
}

// Delete Delete timer from base.
func Delete(tm *Timer) {
	lazyInit()
	base.Delete(tm)
}

// Modify modify timer timeout,and restart immediately.
func Modify(tm *Timer, timeout time.Duration) {
	lazyInit()
	base.Modify(tm, timeout)
}

// AfterFunc like time.AfterFunc
func AfterFunc(d time.Duration, f func()) *Timer { return AddJobFunc(f, d) }
