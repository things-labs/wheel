package timewheel

import (
	"fmt"
	"testing"
	"time"
)

type testJob struct{}

func (sf testJob) Run() {
	fmt.Println("job")
}

func TestDefaultBase(t *testing.T) {
	wl := New(WithGranularity(DefaultGranularity), WithGoPool(nil)).Run()

	defer wl.Close()
	if got := wl.Len(); got != 0 {
		t.Errorf("Len() = %v, want %v", got, 0)
	}
	if got := wl.HasRunning(); got != true {
		t.Errorf("HasRunning() = %v, want %v", got, true)
	}

	e := NewJobFunc(func() {})
	wl.Add(e, time.Millisecond*100)
	wl.Delete(e)
	wl.Modify(e, time.Millisecond*200)
	time.Sleep(time.Second)

	e1 := NewTimer().WithGoroutine()
	wl.Add(e1, time.Millisecond*150)

	e2 := NewTimer().WithGoroutine()
	wl.Add(e2, 0)
	time.Sleep(time.Second)

	// improve couver
	wl.Modify(nil, time.Second)
	wl.Delete(nil)
	wl.Add(nil, time.Second)
}

func ExampleBase_Run() {
	wl := New().Run()

	wl.AddJobFunc(func() {
		fmt.Println("1")
	}, time.Millisecond*100)
	wl.AddJobFunc(func() {
		fmt.Println("2")
	}, time.Millisecond*200)
	wl.AddJob(&testJob{}, time.Millisecond*300)
	time.Sleep(time.Second * 2)
	// Output:
	// 1
	// 2
	// job
}
