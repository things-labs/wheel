package wheel

import (
	"math"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// DefaultGranularity 1ms every 1ms a tick
	DefaultGranularity = time.Millisecond * 1
)

// 主级 + 4个层级共5级 占32位
const (
	tvRBits = 8            // 主级占8位
	tvNBits = 6            // 4个层级各占6位
	tvNNum  = 4            // 层级个数
	tvRSize = 1 << tvRBits // 主级槽个数
	tvNSize = 1 << tvNBits // 每个层级槽个数
	tvRMask = tvRSize - 1  // 主轮掩码
	tvNMask = tvNSize - 1  // 层级掩码
)

// Base time wheel base
type Base struct {
	spokes      []list
	doRunning   *list
	curTick     uint32
	startTime   time.Time
	granularity time.Duration
	rw          sync.RWMutex
	stop        chan struct{}
	running     uint32
}

// Option option for base
type Option func(w *Base)

// WithGranularity override base granularity
func WithGranularity(gra time.Duration) Option {
	return func(w *Base) {
		w.granularity = gra
	}
}

// New new a base wheel
func New(opts ...Option) *Base {
	wl := &Base{
		spokes:      make([]list, tvRSize+tvNSize*tvNNum),
		doRunning:   newList(),
		startTime:   time.Now(),
		granularity: DefaultGranularity,
		stop:        make(chan struct{}),
		curTick:     math.MaxUint32 - 30,
	}
	for _, opt := range opts {
		opt(wl)
	}
	return wl
}

// Run the base in its own goroutine, or no-op if already started.
func (sf *Base) Run() *Base {
	if atomic.CompareAndSwapUint32(&sf.running, 0, 1) {
		go sf.runWork()
	}
	return sf
}

// HasRunning base running status.
func (sf *Base) HasRunning() bool {
	return atomic.LoadUint32(&sf.running) == 1
}

// Close close wait util close
func (sf *Base) Close() error {
	if atomic.CompareAndSwapUint32(&sf.running, 1, 0) {
		sf.stop <- struct{}{}
	}
	return nil
}

// Len the number timer of the base.
func (sf *Base) Len() int {
	var length int

	sf.rw.RLock()
	defer sf.rw.RUnlock()

	for i := 0; i < len(sf.spokes); i++ {
		length += sf.spokes[i].Len()
	}
	length += sf.doRunning.Len()
	return length
}

// AddJob add a job
func (sf *Base) AddJob(job Job, timeout time.Duration) *Timer {
	tm := NewJob(job, timeout)
	sf.Add(tm)
	return tm
}

// AddJobFunc add a job function
func (sf *Base) AddJobFunc(f func(), timeout time.Duration) *Timer {
	return sf.AddJob(JobFunc(f), timeout)
}

// Add add timer to base.
func (sf *Base) Add(tm *Timer, newTimeout ...time.Duration) {
	if tm == nil {
		return
	}
	sf.rw.Lock()
	sf.start(tm, newTimeout...)
	sf.rw.Unlock()
}

// Delete Delete timer from base.
func (sf *Base) Delete(tm *Timer) {
	if tm == nil {
		return
	}
	sf.rw.Lock()
	tm.removeSelf()
	sf.rw.Unlock()
}

// Modify modify timer timeout,and restart immediately.
func (sf *Base) Modify(tm *Timer, timeout time.Duration) {
	if tm == nil {
		return
	}
	sf.rw.Lock()
	sf.start(tm, timeout)
	sf.rw.Unlock()
}

func (sf *Base) runWork() {
	tick := time.NewTimer(sf.granularity)
	for {
		select {
		case now := <-tick.C:
			nano := now.Sub(sf.startTime)
			tick.Reset(nano % sf.granularity)
			sf.rw.Lock()
			for past := uint32(nano/sf.granularity) - sf.curTick; past > 0; past-- {
				sf.curTick++
				index := sf.curTick & tvRMask
				if index == 0 {
					sf.cascade()
				}
				sf.doRunning.SpliceBackList(&sf.spokes[index])
			}

			for sf.doRunning.Len() > 0 {
				tm := sf.doRunning.PopFront()
				sf.rw.Unlock()
				if tm.job != nil {
					if tm.useGoroutine {
						go tm.job.Run()
					} else {
						wrapJob(tm.job)
					}
				}
				sf.rw.Lock()
			}
			sf.rw.Unlock()

		case <-sf.stop:
			tick.Stop()
			return
		}
	}
}

// 层叠计算每一层
func (sf *Base) cascade() {
	for level := 0; ; {
		index := int((sf.curTick >> (uint32)(tvRBits+level*tvNBits)) & tvNMask)
		spoke := sf.spokes[tvRSize+tvNSize*level+index]
		for spoke.Len() > 0 {
			sf.addTimer(spoke.PopFront())
		}
		if level++; !(index == 0 && level < tvNNum) {
			break
		}
	}
}

func (sf *Base) nextTick(next time.Time) uint32 {
	return uint32((next.Sub(sf.startTime) + sf.granularity - 1) / sf.granularity)
}

func (sf *Base) addTimer(tm *Timer) *Timer {
	var spokeIdx int

	next := sf.nextTick(tm.nextTime)
	if idx := next - sf.curTick; idx < tvRSize {
		spokeIdx = int(next & tvRMask)
	} else {
		// 计算在哪一个层级
		level := 0
		for idx >>= tvRBits; idx >= tvNSize && level < (tvNNum-1); level++ {
			idx >>= tvNBits
		}
		spokeIdx = tvRSize + tvNSize*level + int((next>>(uint32)(tvRBits+tvNBits*level))&tvNMask)
	}
	sf.spokes[spokeIdx].PushElementBack(tm)
	return tm
}

func (sf *Base) start(tm *Timer, newTimeout ...time.Duration) {
	tm.removeSelf() // should remove from old list

	if len(newTimeout) > 0 {
		tm.timeout = newTimeout[0]
	}
	tm.nextTime = time.Now().Add(tm.timeout)
	if sf.nextTick(tm.nextTime) == sf.curTick {
		sf.doRunning.PushElementBack(tm)
	} else {
		sf.addTimer(tm)
	}
}
