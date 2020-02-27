package timing

import (
	"time"
)

// Option 选项
type Option func(w *Wheel)

// WithGranularity override interval 时间粒子
func WithGranularity(gra time.Duration) Option {
	return func(w *Wheel) {
		w.granularity = gra
	}
}

// WithInterval override interval 默认条目时间间隔
func WithInterval(interval time.Duration) Option {
	return func(w *Wheel) {
		w.interval = interval
	}
}

// WithGoroutine override hasGoroutine 回调使用goroutine执行
func WithGoroutine(use bool) Option {
	return func(w *Wheel) {
		w.UseGoroutine(use)
	}
}
