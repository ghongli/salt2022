package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	go readLock()
	go writeLock()
	go lock()

	time.Sleep(5 * time.Second)
}

func lock() {
	lock := sync.Mutex{}
	for i := 0; i < 3; i++ {
		lock.Lock()
		// a run-time error if not locked on entry to Unlock.
		lock.Unlock()

		fmt.Println("lock:", i)
	}
}

func readLock() {
	lock := sync.RWMutex{}
	for i := 0; i < 3; i++ {
		lock.RLock()
		// a run-time error if not locked for reading on entry to RUnlock.
		lock.RUnlock()

		fmt.Println("read-lock:", i)
	}
}

func writeLock() {
	lock := sync.RWMutex{}
	for i := 0; i < 3; i++ {
		lock.Lock()
		// a run-time error if not locked for writing on entry to Unlock.
		lock.Unlock()

		fmt.Println("write-lock:", i)
	}
}
