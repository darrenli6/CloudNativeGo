package main

import (
	"fmt"
	"time"
)

func main() {

	stop := make(chan bool)

	go func() {

		for {

			select {
			case <-stop:
				fmt.Println("stop channel")
				return
			default:
				fmt.Println("Printing")
				time.Sleep(time.Second)

			}
		}
	}()

	time.Sleep(5 * time.Second)
	fmt.Println("stoppping")
	stop <- true
	time.Sleep(5 * time.Second)
}

/*
select多路复用select 就是监听 IO 操作，当 IO 操作发生时，触发相应的动作。 select的用法与switch非常类似，由select开始一个新的选择块，每个选择条件由case语句来描述。与switch语句可以选择任何可使用相等比较的条件相比，select有比较多的限制，其中最大的一条限制就是每个case语句里必须是一个IO操作，确切的说，应该是一个面向channel的IO操作。可处理一个或多个channel的发送/接收操作。如果多个case同时满足，select会随机选择一个。对于没有case的select{}会一直等待，可用于阻塞main函数。
*/
