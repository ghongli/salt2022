package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
)

// go tool pprof -http :8081 /tmp/profile/cpu
var cpuprofile = flag.String("cpuprofile", "/tmp/profile/cpu", "write cpu profile to file")

func main() {
	flag.Parse()
	f, err := os.Create(*cpuprofile)
	if err != nil {
		log.Fatal(err)
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	var result int
	for i := 0; i < 100000000; i++ {
		result += i
	}
	log.Println("result: ", result)
}
