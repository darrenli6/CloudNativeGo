package main

import (
	"fmt"
	"sync"
	"time"
)

/*
当你需要将一个功能拆成不同的job去执行，然后等全部job完成了，再继续执行下面的程序，那你就需要waitgroup
那这有一个问题，如果多个job，其中某个job卡住了，长期不返回结果了，怎么中断他？我们可以使用select+channel
*/
func main() {

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("thread 1 is completed")
		wg.Done()
	}()

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("thread 2 is completed")
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("all thread are completed ")

}
