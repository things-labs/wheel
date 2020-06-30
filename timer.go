package wheel

import (
	"time"
)

// Timer which hold the timer instance
type Timer struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Timer

	// The list to which this element belongs.
	list *list

	// nextTime 下一次运行时间  0: 表示未运行,或未启动
	nextTime time.Time
	// timeout 超时时间
	timeout time.Duration
	// 任务
	job Job
	// useGoroutine
	useGoroutine bool
}

// NewTimer new a timer with a empty job,
func NewTimer(timeout time.Duration) *Timer {
	return &Timer{
		timeout: timeout,
		job:     emptyJob{},
	}
}

// NewJob 新建一个条目,条目未启动定时
func NewJob(job Job, timeout time.Duration) *Timer {
	return NewTimer(timeout).WithJob(job)
}

// NewJobFunc 新建一个条目,条目未启动定时
func NewJobFunc(f func(), timeout time.Duration) *Timer {
	return NewTimer(timeout).WithJobFunc(f)
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
