package main

import (
	"log"
	"time"

	timewheel "github.com/things-go/timewheel"
)

func main() {
	tm := timewheel.NewTimer()
	tm.WithJobFunc(func() {
		log.Println("hello world")
		timewheel.Add(tm, time.Second)
	})
	timewheel.Add(tm, time.Second)
	time.Sleep(time.Second * 60)
}
