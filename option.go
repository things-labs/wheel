package timing

import (
	"time"
)

// Option 选项
type Option func(w *Wheel)

// WithGranularity override timeout 时间粒子
func WithGranularity(gra time.Duration) Option {
	return func(w *Wheel) {
		w.granularity = gra
	}
}

// WithGoroutine override useGoroutine 回调使用goroutine执行
func WithGoroutine(use bool) Option {
	return func(w *Wheel) {
		w.UseGoroutine(use)
	}
}
