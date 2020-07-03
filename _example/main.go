package main

import (
	"log"
	"time"

	"github.com/thinkgos/wheel"
)

func main() {
	base := wheel.New().Run()

	tm := wheel.NewTimer()
	tm.WithJobFunc(func() {
		log.Println("hello world")
		base.Add(tm, time.Second)
	})
	base.Add(tm, time.Second)
	time.Sleep(time.Second * 60)
}
