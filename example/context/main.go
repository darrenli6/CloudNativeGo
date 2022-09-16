package main

import (
	"context"
	"fmt"
	"time"
)

// web初始化一些值
// 传递到全局变量去

func main() {

	baseCtx := context.Background()

	ctx := context.WithValue(baseCtx, "a", "b")

	go func(c context.Context) {
		fmt.Println(c.Value("a"))
	}(ctx)

	timeOutCtx, cancel := context.WithTimeout(baseCtx, time.Second)

	defer cancel()
	// 处理一些时间比较长 不想等待

	go func(c context.Context) {
		ticker := time.NewTicker(time.Second)
		for _ = range ticker.C {
			select {
			case <-c.Done():
				fmt.Println("child process interupt")
				return
			default:
				fmt.Println("enter default")

			}
		}
	}(timeOutCtx)

	select {
	case <-timeOutCtx.Done():
		time.Sleep(time.Second)
		fmt.Println("main process exit")
	}

}
