package main

import (
	"fmt"
	"sync"
	"time"
)

type Queue struct {
	queue []string
	cond  *sync.Cond
}

// 设计生产者和消费者

func main() {
	q := Queue{
		queue: []string{},
		cond:  sync.NewCond(&sync.Mutex{}),
	}

	// 生产者
	go func() {
		for {
			q.Enqueue("a")
			time.Sleep(2 * time.Second)
		}
	}()

	// 消费者

	for {
		q.Dequque()
		time.Sleep(1 * time.Second)

	}

}

func (q *Queue) Enqueue(item string) {

	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	// 往队列加入东西
	q.queue = append(q.queue, item)
	fmt.Printf("putting %s to queue  notify all \n", item)
	// 通知消费者
	q.cond.Broadcast()
}

func (q *Queue) Dequque() string {
	q.cond.L.Lock()

	defer q.cond.L.Unlock()

	// 如果队列为空 怎么办 ？
	if len(q.queue) == 0 {
		fmt.Println("no data avaiable ")
		// 等待
		q.cond.Wait()
	}
	// 取东西
	result := q.queue[0]
	// 切片
	q.queue = q.queue[1:]

	return result

}
