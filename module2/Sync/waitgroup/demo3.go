package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var fmeng = false

func producer(threadID int, wg *sync.WaitGroup, ch chan string) {

	count := 0

	for !fmeng {

		time.Sleep(time.Second)
		count++
		data := strconv.Itoa(threadID) + "-----" + strconv.Itoa(count)
		fmt.Printf("producer %s \n", data)
		ch <- data
	}
	wg.Done()
}

func consumer(wg *sync.WaitGroup, ch chan string) {
	for data := range ch {
		time.Sleep(time.Second)
		fmt.Printf("consumer %s \n", data)
	}
	wg.Done()
}

func main() {

	//三个生产者和二个消费者
	chanStream := make(chan string, 10)

	// 生产者和消费者计数器
	wgPd := new(sync.WaitGroup)
	wgCs := new(sync.WaitGroup)

	for i := 0; i < 3; i++ {
		wgPd.Add(1)
		go producer(i, wgPd, chanStream)
	}

	for i := 0; i < 2; i++ {
		wgCs.Add(1)
		go consumer(wgCs, chanStream)
	}

	//制造超时
	go func() {
		time.Sleep(time.Second * 10)
		fmeng = true
	}()
	wgPd.Wait()

	close(chanStream)

	wgCs.Wait()
}
