package wheel

import (
	"time"

	"github.com/thinkgos/list"
)

// entry 条目
type entry struct {
	// next 下一次运行时间  0: 表示未运行,或未启动
	next time.Time
	// timeout 超时时间
	timeout time.Duration
	// 任务
	job Job
	// useGoroutine
	useGoroutine bool
}

// Timer which hold the timer instance
type Timer list.Element

// NewTimer new a timer with a empty job,
func NewTimer(timeout time.Duration) *Timer {
	return &Timer{
		Value: &entry{
			timeout: timeout,
			job:     emptyJob{},
		},
	}
}

func (sf *Timer) WithGoroutine() *Timer {
	sf.getEntry().useGoroutine = true
	return sf
}

func (sf *Timer) WithTimeout(timeout time.Duration) *Timer {
	sf.getEntry().timeout = timeout
	return sf
}

func (sf *Timer) WithJob(job Job) *Timer {
	sf.getEntry().job = job
	return sf
}

func (sf *Timer) WithJobFunc(f func()) *Timer {
	return sf.WithJob(JobFunc(f))
}

func (sf *Timer) getEntry() *entry {
	return sf.Value.(*entry)
}
