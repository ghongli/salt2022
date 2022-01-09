package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("recover err: %v\n", err)
		}
	}()
	
	// unsafeWrite()
	safeWrite()

	time.Sleep(time.Second)
}

func unsafeWrite() {
	m := map[int]int{}
	for i := 0; i < 100; i++ {
		i := i
		go func() {
			// runtime.gopanic() panic 只能触发当前协程的 defer func
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("unsafe writer's recover err: %v\n", err)
				}
			}()
			
			// runtime.goexit() 会产生一个 _panic 实例，相应也会处理 defer 函数，但不能被 recover 函数恢复
			m[i] = i
		}()
	}
}

func safeWrite() {
	m := SafeMap{
		safeMap: map[int]int{},
		Mutex:   sync.Mutex{},
	}

	for i := 0; i < 100; i++ {
		i := i
		go func() {
			m.Write(i, i)
		}()
	}
}

type (
	SafeMap struct {
		safeMap map[int]int
		sync.Mutex
	}
)

func (m *SafeMap) Read(k int) (int, bool) {
	m.Lock()
	defer m.Unlock()

	result, ok := m.safeMap[k]
	return result, ok
}

func (m *SafeMap) Write(k, v int) {
	m.Lock()
	defer m.Unlock()

	m.safeMap[k] = v
}
