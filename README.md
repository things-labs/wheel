# wheel timer

golang time wheel library, which similar linux time wheel

[![GoDoc](https://godoc.org/github.com/things-labs/wheel?status.svg)](https://godoc.org/github.com/things-labs/wheel)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/things-labs/wheel?tab=doc)
[![Build Status](https://travis-ci.com/things-labs/wheel.svg)](https://travis-ci.com/things-labs/wheel)
[![codecov](https://codecov.io/gh/things-labs/wheel/branch/master/graph/badge.svg)](https://codecov.io/gh/things-labs/wheel)
![Action Status](https://github.com/things-labs/wheel/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/things-labs/wheel)](https://goreportcard.com/report/github.com/things-labs/wheel)
[![Licence](https://img.shields.io/github/license/things-labs/wheel)](https://raw.githubusercontent.com/things-labs/wheel/master/LICENSE)
[![Tag](https://img.shields.io/github/v/tag/things-labs/wheel)](https://github.com/things-labs/wheel/tags)

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
    go get github.com/things-labs/wheel
```

Then import the wheel package into your own code.
```bash
    import "github.com/things-labs/wheel"
```

### Example

---

[embedmd]# (_example/main.go go)
```go
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
```

### References

---
