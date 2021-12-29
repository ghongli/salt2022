package main

import (
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(1)
	
	go func() {
		// 类 for {} 无限循环 <- 无函数调用的无限循环
		println("for {}")
		select {} // 永远阻塞
	}()
	
	time.Sleep(1 * time.Second) // 系统调用，出让执行权给上面的协程
	println("Done")
}
