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

// Len 条目个数
func Len() int {
	lazyInit()
	return base.Len()
}

// AddJob 添加任务
func AddJob(job Job, timeout time.Duration) *Timer {
	lazyInit()
	return base.AddJob(job, timeout)
}

// AddJobFunc 添加任务函数
func AddJobFunc(f func(), timeout time.Duration) *Timer {
	lazyInit()
	return AddJob(JobFunc(f), timeout)
}

// Add 启动或重始启动e的计时
func Add(tm *Timer, newTimeout ...time.Duration) {
	lazyInit()
	base.Add(tm, newTimeout...)
}

// Delete 删除条目
func Delete(tm *Timer) {
	lazyInit()
	base.Delete(tm)
}

// Modify 修改条目的周期时间,重置计数且重新启动定时器
func Modify(tm *Timer, timeout time.Duration) {
	lazyInit()
	base.Modify(tm, timeout)
}
