package main

import (
	"fmt"
	"runtime"
	"time"
)

const NUM = 100

func main() {

	runtime.GOMAXPROCS(2)
	go a()
	go b()

	time.Sleep(time.Hour)
}

func a() {
	for i := 0; i < NUM; i++ {
		fmt.Println("A:", i)
	}
}

func b() {
	for i := 0; i < NUM; i++ {
		fmt.Println("B:", i)
	}
}
