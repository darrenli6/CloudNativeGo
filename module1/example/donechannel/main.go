package main

import (
	"fmt"
	"time"
)

func main() {
	message := make(chan int, 10)

	done := make(chan bool)

	defer close(message)
	//consumer
	go func() {

		timer := time.NewTicker(time.Second)
		for _ = range timer.C {
			select {
			case <-done:
				fmt.Println("child process is interrupted")
				break
			default:
				fmt.Printf("send message %d \n", <-message)

			}
		}

	}()

	// producer
	for i := 0; i < 10; i++ {
		fmt.Printf("put data  %d \n", i)
		message <- i

	}

	time.Sleep(5 * time.Second)
	close(done)
	time.Sleep(time.Second)
	fmt.Println("main process exit")

}
