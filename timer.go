package wheel

import (
	"time"
)

// Timer consists of a schedule and the func to execute on that schedule.
type Timer struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Timer

	// The list to which this element belongs.
	list *list

	// next time the job will run, or the zero time if Base has not been
	// started or this entry is unsatisfiable
	nextTime time.Time
	// job is the thing that want to run.
	job Job
	// use goroutine
	useGoroutine bool
}

// NewTimer new a timer with a empty job,
func NewTimer() *Timer {
	return &Timer{
		job: emptyJob{},
	}
}

// NewJob new timer with job.
func NewJob(job Job) *Timer {
	return NewTimer().WithJob(job)
}

// NewJobFunc new timer with job function.
func NewJobFunc(f func()) *Timer {
	return NewTimer().WithJobFunc(f)
}

// WithGoroutine with goroutine
func (sf *Timer) WithGoroutine() *Timer {
	sf.useGoroutine = true
	return sf
}

// WithJob with job.
func (sf *Timer) WithJob(job Job) *Timer {
	sf.job = job
	return sf
}

// WithJobFunc with job function
func (sf *Timer) WithJobFunc(f func()) *Timer {
	return sf.WithJob(JobFunc(f))
}
