package main

import (
	"log"
	"time"

	"github.com/thinkgos/wheel"
)

func main() {
	base := wheel.New()
	base.Run()

	tm := wheel.NewTimer(time.Second)
	tm.WithJobFunc(func() {
		log.Println("hello world")
		base.Add(tm)
	})
	base.Add(tm)
	time.Sleep(time.Second * 60)
}
