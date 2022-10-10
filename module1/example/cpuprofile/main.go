package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "/tmp/cpuprofile", "write cpu profile ")

func main() {

	flag.Parse()

	f, err := os.Create(*cpuprofile)
	if err != nil {
		log.Panic(err)
	}

	pprof.StartCPUProfile(f)

	defer pprof.StopCPUProfile()

	var result int
	for i := 0; i < 10000000; i++ {
		result += i
	}

	fmt.Printf("result is %d \n", result)

}
