package main

//#cgo LDFLAGS:
//char* allocMemory();
import "C"
import (
	"fmt"
	"time"
)

func main() {
	// only loop 5 times to avoid exhausting the host memory
	holder := []*C.char{}
	for i := 1; i < 5; i++ {
		fmt.Printf("Allocating %dMb memory, raw memory is %d\n", i*100, i*100*1024*1025)
		// hold the memory, otherwise will be forced by GC
		holder = append(holder, (*C.char)(C.allocMemory()))
		time.Sleep(time.Minute)
	}
}
