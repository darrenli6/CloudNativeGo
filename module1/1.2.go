package main

import (
	"fmt"
	"time"
)

const NUM = 10

func main() {

	message := make(chan int, NUM)
	defer close(message)

	go func() {
		// consumer
		for i := 0; i < NUM; i++ {

			fmt.Printf("send message %d \n", <-message)

			time.Sleep(time.Second)
		}

	}()

	// product
	for i := 0; i < NUM; i++ {
		fmt.Printf("put  %d  to messages \n", i)
		message <- i

		time.Sleep(time.Second)
	}

	time.Sleep(time.Hour)

}
