package main

import (
	"log"
	"time"

	"github.com/things-labs/wheel"
)

func main() {
	tm := wheel.NewTimer()
	tm.WithJobFunc(func() {
		log.Println("hello world")
		wheel.Add(tm, time.Second)
	})
	wheel.Add(tm, time.Second)
	time.Sleep(time.Second * 60)
}
