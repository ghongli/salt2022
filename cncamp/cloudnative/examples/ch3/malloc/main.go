package main

/**
CGO_ENABLED=1，进行编译时，会将文件中引用libc的库(如常用的net包)，以动态链接的方式生成目标文件 -- 默认；
CGO_ENABLED=0，进行编译时，则会把在目标文件中未定义的符号(外部函数)一起链接到可执行文件中；
ldd target_file 查看引用的动态链接库
nm  target_file 查看符号表
nm main0 | grep ' U '

ldd: print shared library depencies
nm:  list symbols from object files
*/

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
