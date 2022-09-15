package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// waitBySleep()
	// waitByChannel()
	waitByWg()
}
func waitBySleep() {
	for i := 0; i < 100; i++ {
		go fmt.Println(i)
	}
	time.Sleep(time.Second)
}

func waitByChannel() {
	c := make(chan bool, 100)

	for i := 0; i < 100; i++ {
		go func(i int) {
			fmt.Println(i)
			c <- true

		}(i)
	}

	// 这种方式控制主线程什么时候退出
	for i := 0; i < 100; i++ {
		<-c
	}

}

func waitByWg() {
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			fmt.Println(i)
			wg.Done()
		}(i)
	}
	// 所有的wg.done之后
	wg.Wait()
}
