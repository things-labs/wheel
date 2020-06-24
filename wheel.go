package timing

import (
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/thinkgos/list"
)

// num 定义
const (
	OneShot = 1
	Persist = 0
)

const (
	// DefaultGranularity 默认时基精度,意思是每xx时间一个tick
	DefaultGranularity = time.Millisecond * 1
)
const (
	// 主级 + 4个层级共5级 占32位
	tvRBits = 8            // 主级占8位
	tvNBits = 6            // 4个层级各占6位
	tvNNum  = 4            // 层级个数
	tvRSize = 1 << tvRBits // 主级槽个数
	tvNSize = 1 << tvNBits // 每个层级槽个数
	tvRMask = tvRSize - 1  // 主轮掩码
	tvNMask = tvNSize - 1  // 层级掩码
)

// Entry 条目
type Entry struct {
	// next 下一次运行时间  0: 表示未运行,或未启动
	next time.Time
	// 任务已经执行的次数
	count uint32
	// 任务需要执行的次数
	number uint32
	// 时间间隔
	interval time.Duration
	// 任务
	job Job
}

// 内部使用,条目
func entry(e *list.Element) *Entry {
	return e.Value.(*Entry)
}

// Wheel 时间轮实现
type Wheel struct {
	spokes       []*list.List // 轮的槽
	doNow        *list.List
	curTick      uint32
	startTime    time.Time
	granularity  time.Duration
	rw           sync.RWMutex
	stop         chan struct{}
	running      uint32
	hasGoroutine uint32
}

// Timer which hold the timer instance
type Timer *list.Element

// New new a wheel
func New(opts ...Option) *Wheel {
	wl := &Wheel{
		spokes:      make([]*list.List, tvRSize+tvNSize*tvNNum),
		doNow:       list.New(),
		startTime:   time.Now(),
		granularity: DefaultGranularity,
		stop:        make(chan struct{}),
	}

	wl.curTick = math.MaxUint32 - 30
	for i := 0; i < len(wl.spokes); i++ {
		wl.spokes[i] = list.New()
	}

	for _, opt := range opts {
		opt(wl)
	}

	return wl
}

// UseGoroutine use goroutine or callback
func (sf *Wheel) UseGoroutine(use bool) {
	var val uint32
	if use {
		val = 1
	}
	atomic.StoreUint32(&sf.hasGoroutine, val)
}

// Run 运行,不阻塞
func (sf *Wheel) Run() *Wheel {
	if atomic.CompareAndSwapUint32(&sf.running, 0, 1) {
		go sf.runWork()
	}

	return sf
}

// HasRunning 运行状态
func (sf *Wheel) HasRunning() bool {
	return atomic.LoadUint32(&sf.running) == 1
}

// Close close wait util close
func (sf *Wheel) Close() error {
	if atomic.CompareAndSwapUint32(&sf.running, 1, 0) {
		sf.stop <- struct{}{}
	}
	return nil
}

// Len 条目个数
func (sf *Wheel) Len() int {
	var length int

	sf.rw.RLock()
	for i := 0; i < len(sf.spokes); i++ {
		length += sf.spokes[i].Len()
	}
	length += sf.doNow.Len()
	sf.rw.RUnlock()
	return length
}

func (sf *Wheel) nextTick(next time.Time) uint32 {
	return uint32((next.Sub(sf.startTime) + sf.granularity - 1) / sf.granularity)
}

// NewTimer new a timer which mount a empty job, 条目未启动
func (sf *Wheel) NewTimer(num uint32, interval time.Duration) Timer {
	return &list.Element{
		Value: &Entry{
			number:   num,
			interval: interval,
			job:      JobFunc(func() {}),
		},
	}
}

// MountJobOnTimer mount a job on timer
func (sf *Wheel) MountJobOnTimer(tm Timer, job Job) Timer {
	entry(tm).job = job
	return tm
}

// MountJobFuncOnTimer mount a job function on timer
func (sf *Wheel) MountJobFuncOnTimer(tm Timer, f JobFunc) Timer {
	return sf.MountJobOnTimer(tm, f)
}

// NewJob 新建一个条目,条目未启动定时
func (sf *Wheel) NewJob(job Job, num uint32, interval time.Duration) Timer {
	return sf.MountJobOnTimer(sf.NewTimer(num, interval), job)
}

// NewJobFunc 新建一个条目,条目未启动定时
func (sf *Wheel) NewJobFunc(f JobFunc, num uint32, interval time.Duration) Timer {
	return sf.MountJobFuncOnTimer(sf.NewTimer(num, interval), f)
}

// AddJob 添加任务
func (sf *Wheel) AddJob(job Job, num uint32, interval time.Duration) Timer {
	e := sf.NewJob(job, num, interval)
	entry := entry(e)
	entry.next = time.Now().Add(entry.interval)

	sf.rw.Lock()
	defer sf.rw.Unlock()

	if sf.nextTick(entry.next) == sf.curTick {
		return sf.doNow.PushElementBack(e)
	}
	return sf.addTimer(e)
}

// AddOneShotJob 添加一次性任务
func (sf *Wheel) AddOneShotJob(job Job, interval time.Duration) Timer {
	return sf.AddJob(job, OneShot, interval)
}

// AddPersistJob 添加周期性任务
func (sf *Wheel) AddPersistJob(job Job, interval time.Duration) Timer {
	return sf.AddJob(job, Persist, interval)
}

// AddJobFunc 添加任务函数
func (sf *Wheel) AddJobFunc(f JobFunc, num uint32, interval time.Duration) Timer {
	return sf.AddJob(f, num, interval)
}

// AddOneShotJobFunc 添加一次性任务函数
func (sf *Wheel) AddOneShotJobFunc(f JobFunc, interval time.Duration) Timer {
	return sf.AddJob(f, OneShot, interval)
}

// AddPersistJobFunc 添加周期性函数
func (sf *Wheel) AddPersistJobFunc(f JobFunc, interval time.Duration) Timer {
	return sf.AddJob(f, Persist, interval)
}

func (sf *Wheel) start(e *list.Element, newTimeout ...time.Duration) *Wheel {
	e.RemoveSelf() // should remove from old list
	entry := entry(e)
	entry.count = 0
	entry.next = time.Now().Add(append(newTimeout, entry.interval)[0])

	sf.addTimer(e)

	return sf
}

// Start 启动或重始启动e的计时
func (sf *Wheel) Start(tm Timer, newTimeout ...time.Duration) *Wheel {
	if tm == nil {
		return sf
	}

	sf.rw.Lock()
	sf.start(tm, newTimeout...)
	sf.rw.Unlock()

	return sf
}

// Delete 删除条目
func (sf *Wheel) Delete(tm Timer) *Wheel {
	if tm == nil {
		return sf
	}

	sf.rw.Lock()
	(*list.Element)(tm).RemoveSelf()
	sf.rw.Unlock()

	return sf
}

// Modify 修改条目的周期时间,重置计数且重新启动定时器
func (sf *Wheel) Modify(tm Timer, interval time.Duration) *Wheel {
	if tm == nil {
		return sf
	}

	sf.rw.Lock()
	entry(tm).interval = interval
	sf.start(tm)
	sf.rw.Unlock()

	return sf
}

func (sf *Wheel) runWork() {
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
				sf.doNow.SpliceBackList(sf.spokes[index])
			}

			for sf.doNow.Len() > 0 {
				e := sf.doNow.PopFront()
				entry := entry(e)

				entry.count++
				if entry.number == 0 || entry.count < entry.number {
					entry.next = now.Add(entry.interval)
					sf.addTimer(e)
				}
				sf.rw.Unlock()
				if atomic.LoadUint32(&sf.hasGoroutine) == 1 {
					go entry.job.Run()
				} else {
					wrapJob(entry.job)
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
func (sf *Wheel) cascade() {
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

func (sf *Wheel) addTimer(tm Timer) *list.Element {
	var spokeIdx int

	next := sf.nextTick(entry(tm).next)
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

	return sf.spokes[spokeIdx].PushElementBack(tm)
}
