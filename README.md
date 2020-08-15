# wheel timer

golang time wheel library, which similar linux time wheel

[![GoDoc](https://godoc.org/github.com/thinkgos/wheel?status.svg)](https://godoc.org/github.com/thinkgos/wheel)
[![Build Status](https://travis-ci.org/thinkgos/wheel.svg?branch=master)](https://travis-ci.org/thinkgos/wheel)
[![codecov](https://codecov.io/gh/thinkgos/wheel/branch/master/graph/badge.svg)](https://codecov.io/gh/thinkgos/wheel)
![Action Status](https://github.com/thinkgos/wheel/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/thinkgos/wheel)](https://goreportcard.com/report/github.com/thinkgos/wheel)
[![Licence](https://img.shields.io/github/license/thinkgos/wheel)](https://raw.githubusercontent.com/thinkgos/wheel/master/LICENSE)  

### Feature

 - Five-level time wheel: main level and four levels.
 - insert,delete,modify,scan item time complexity O(1).
 - the default time granularity is 1ms.
 - The maximum time is limited by the accuracy of the time base. The time granularity is 1ms, 
 and the maximum time can be 49.71 days. so the maximum time is 49.71 days * (granularity/1ms)
 - There is the internal wheel timer base with granularity 1ms,it lazies init internal until you first used.
 - **NOTE:do not use Time consuming task @ timer callback function,you can with `WithGoroutine`** 


### Installation

Use go get.
```bash
    go get github.com/thinkgos/wheel
```

Then import the wheel package into your own code.
```bash
    import "github.com/thinkgos/wheel"
```

### Example

---

```go
import (
	"log"
	"time"

	"github.com/thinkgos/wheel"
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
```

### References

---
