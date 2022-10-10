package main

import (
	"fmt"
	"sync"
	"time"
)

// 等待一组goroutine 返回

func main() {
	// waitBySleep()
	waitByChannel()
	// waitByWg()
}
func waitBySleep() {
	for i := 0; i < 100; i++ {
		go fmt.Println(i)
	}
	// sleep足够长的时间工作就能结束
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
	// 如果不够100 就阻塞   直到退出
	// 上下要协调好
	for i := 0; i < 100; i++ {
		<-c
	}

}

func waitByWg() {
	wg := sync.WaitGroup{}
	// 加入100个线程
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			fmt.Println(i)
			// 如果不够100  还是阻塞在这里
			wg.Done()
		}(i)
	}
	// 针对多个线程协调
	// 所有的wg.done之后
	//如果不够100  还是阻塞在这里
	wg.Wait()
}
