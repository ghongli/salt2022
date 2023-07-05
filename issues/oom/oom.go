package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func main() {
	for i := 0; i < 4; i++ {
		queryAll()
		fmt.Printf("goroutines: %d\n", runtime.NumGoroutine())
	}
}

func queryAll() int {
	ch := make(chan int)
	for i := 0; i < 3; i++ {
		go func() {
			ch <- query()
		}()
	}

	return <-ch
}

func query() int {
	n := rand.Intn(100)
	time.Sleep(time.Duration(n) * time.Millisecond)
	return n
}
