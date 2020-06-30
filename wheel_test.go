package wheel

import (
	"fmt"
	"testing"
	"time"
)

func TestWheel(t *testing.T) {
	if got := Len(); got != 0 {
		t.Errorf("Len() = %v, want %v", got, 0)
	}

	e := NewJobFunc(func() {}, time.Millisecond*100)
	Add(e)
	Delete(e)
	Modify(e, time.Millisecond*200)
	time.Sleep(time.Second)

	e1 := NewTimer(time.Millisecond * 100).WithGoroutine()
	Add(e1, time.Millisecond*150)

	e2 := NewTimer(time.Millisecond * 100).WithGoroutine()
	Add(e2, 0)
	time.Sleep(time.Second)

	// improve couver
	Modify(nil, time.Second)
	Delete(nil)
	Add(nil)
}

func ExampleBase_Len() {
	AddJobFunc(func() {
		fmt.Println("1")
	}, time.Millisecond*100)
	AddJobFunc(func() {
		fmt.Println("2")
	}, time.Millisecond*200)
	AddJob(&testJob{}, time.Millisecond*300)
	time.Sleep(time.Second * 2)
	// Output:
	// 1
	// 2
	// job
}
