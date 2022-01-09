package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// waitByWG()
	waitByChannel()
}

func waitByWG() {
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		i := i
		go func() {
			fmt.Println(i)
			wg.Done()
		}()
	}
	wg.Wait()
}

func waitByChannel() {
	ch := make(chan int, 100)
	for i := 0; i < 100; i++ {
		i := i
		go func() {
			fmt.Println(i)
			ch <- i
		}()
	}

	for i := 0; i < 100; i++ {
		<-ch
	}
}

func waitBySleep() {
	for i := 0; i < 100; i++ {
		go fmt.Println(i)
	}

	time.Sleep(time.Second)
}
