package wheel

import (
	"sync"
	"time"
)

var once sync.Once
var base *Base

func lazyInit() {
	once.Do(func() {
		base = New().Run()
	})
}

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
func Add(tm *Timer, newTimeout ...time.Duration) {
	lazyInit()
	base.Add(tm, newTimeout...)
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
