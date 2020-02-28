package timing

import (
	"fmt"
	"testing"
	"time"
)

type emptyJob struct{}

func (emptyJob) Run() {}

type testJob struct{}

func (sf testJob) Run() {
	fmt.Println("job")
}

func TestDefaultTiming(t *testing.T) {
	wl := New(WithGranularity(DefaultGranularity),
		WithGoroutine(true)).Run()
	wl.UseGoroutine(false)
	defer wl.Close()
	if got := wl.Len(); got != 0 {
		t.Errorf("Len() = %v, want %v", got, 0)
	}
	if got := wl.HasRunning(); got != true {
		t.Errorf("HasRunning() = %v, want %v", got, true)
	}

	wl.AddPersistJobFunc(func() {}, time.Millisecond*100)
	wl.AddPersistJob(&emptyJob{}, time.Second)
	e := wl.NewJobFunc(func() {}, 2, time.Millisecond*100)
	wl.Start(e)
	wl.Delete(e)
	wl.Modify(e, time.Millisecond*200)
	time.Sleep(time.Second)
}

func ExampleWheel_Run() {
	wl := New().Run()

	wl.AddOneShotJobFunc(func() {
		fmt.Println("1")
	}, time.Millisecond*100)
	wl.AddJobFunc(func() {
		fmt.Println("2")
	}, OneShot, time.Millisecond*200)
	wl.AddOneShotJob(&testJob{}, time.Millisecond*300)
	wl.AddJob(&testJob{}, 2, time.Millisecond*400)
	wl.UseGoroutine(true)
	time.Sleep(time.Second * 2)
	// Output:
	// 1
	// 2
	// job
	// job
	// job
}
