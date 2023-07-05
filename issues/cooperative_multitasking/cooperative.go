package main

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"
)

func main() {
	// a single thread(think GOMAXPROCS=1)
	var ops uint64 = 0
	for i := 0; i < 50; i++ {
		go func() {
			for {
				atomic.AddUint64(&ops, 1)
				// workaround:
				runtime.Gosched()
			}
		}()

		fmt.Printf("goroutines: %d\n", runtime.NumGoroutine())
	}

	time.Sleep(time.Second)

	opsFinal := atomic.LoadUint64(&ops)
	fmt.Println("ops:", opsFinal)
}
