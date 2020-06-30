package wheel

import (
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/thinkgos/list"
)

const (
	// DefaultGranularity 默认时基精度,意思是每xx时间一个tick
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

// Base 时间轮实现
type Base struct {
	spokes      []list.List // 轮的槽
	doRunning   *list.List
	curTick     uint32
	startTime   time.Time
	granularity time.Duration
	rw          sync.RWMutex
	stop        chan struct{}
	running     uint32
}

// Option 选项
type Option func(w *Base)

// WithGranularity override timeout 时间粒子
func WithGranularity(gra time.Duration) Option {
	return func(w *Base) {
		w.granularity = gra
	}
}

// New new a wheel
func New(opts ...Option) *Base {
	wl := &Base{
		spokes:      make([]list.List, tvRSize+tvNSize*tvNNum),
		doRunning:   list.New(),
		startTime:   time.Now(),
		granularity: DefaultGranularity,
		stop:        make(chan struct{}),
		curTick:     math.MaxUint32 - 30,
	}
	for i := 0; i < len(wl.spokes); i++ {
		wl.spokes[i].Init()
	}
	for _, opt := range opts {
		opt(wl)
	}
	return wl
}

// Run 运行,不阻塞
func (sf *Base) Run() *Base {
	if atomic.CompareAndSwapUint32(&sf.running, 0, 1) {
		go sf.runWork()
	}
	return sf
}

// HasRunning 运行状态
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

// Len 条目个数
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

// NewJob 新建一个条目,条目未启动定时
func NewJob(job Job, timeout time.Duration) *Timer {
	return NewTimer(timeout).WithJob(job)
}

// NewJobFunc 新建一个条目,条目未启动定时
func NewJobFunc(f func(), timeout time.Duration) *Timer {
	return NewTimer(timeout).WithJob(JobFunc(f))
}

// AddJob 添加任务
func (sf *Base) AddJob(job Job, timeout time.Duration) *Timer {
	tm := NewJob(job, timeout)
	e := tm.getEntry()
	e.next = time.Now().Add(e.timeout)

	sf.rw.Lock()
	defer sf.rw.Unlock()

	if sf.nextTick(e.next) == sf.curTick {
		sf.doRunning.PushElementBack((*list.Element)(tm))
		return tm
	}
	return sf.addTimer(tm)
}

// AddJobFunc 添加任务函数
func (sf *Base) AddJobFunc(f func(), interval time.Duration) *Timer {
	return sf.AddJob(JobFunc(f), interval)
}

// Add 启动或重始启动e的计时
func (sf *Base) Add(tm *Timer, newTimeout ...time.Duration) *Base {
	if tm == nil {
		return sf
	}

	sf.rw.Lock()
	defer sf.rw.Unlock()

	return sf.start(tm, newTimeout...)
}

// Delete 删除条目
func (sf *Base) Delete(tm *Timer) *Base {
	if tm == nil {
		return sf
	}

	sf.rw.Lock()
	(*list.Element)(tm).RemoveSelf()
	sf.rw.Unlock()

	return sf
}

// Modify 修改条目的周期时间,重置计数且重新启动定时器
func (sf *Base) Modify(tm *Timer, timeout time.Duration) *Base {
	if tm == nil {
		return sf
	}

	sf.rw.Lock()
	defer sf.rw.Unlock()

	tm.getEntry().timeout = timeout
	return sf.start(tm)
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
				tm := (*Timer)(sf.doRunning.PopFront())
				e := tm.getEntry()

				sf.rw.Unlock()
				if e.useGoroutine {
					go e.job.Run()
				} else {
					wrapJob(e.job)
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
			sf.addTimer((*Timer)(spoke.PopFront()))
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

	next := sf.nextTick(tm.getEntry().next)
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
	sf.spokes[spokeIdx].PushElementBack((*list.Element)(tm))
	return tm
}

func (sf *Base) start(tm *Timer, newTimeout ...time.Duration) *Base {
	(*list.Element)(tm).RemoveSelf() // should remove from old list
	e := tm.getEntry()

	timeout := e.timeout
	if len(newTimeout) > 0 {
		timeout = newTimeout[0]
	}
	e.next = time.Now().Add(timeout)

	sf.addTimer(tm)
	return sf
}
