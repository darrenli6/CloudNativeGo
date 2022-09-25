package main

import (
	"fmt"
	"sync"
	"time"
)

// 不限制并发读，
// 只限制并发写和并发读写
func main() {

	// go rlock()
	go wlock()
	// go lock()

	time.Sleep(5 * time.Second)
}

func lock() {
	// 这个同样会互斥
	// 标准锁
	lock := sync.Mutex{}

	for i := 0; i < 3; i++ {
		// 第二次获取不了锁 卡死
		lock.Lock()
		defer lock.Unlock()
		fmt.Println("lock ", i)

	}
}

func rlock() {
	// 读写分离锁
	lock := sync.RWMutex{}

	for i := 0; i < 3; i++ {

		// 读锁
		lock.RLock()
		// 读锁不互斥，所以会循环三次，然后将锁释放掉
		defer lock.RUnlock()
		fmt.Println("rlock ", i)

	}
}

func wlock() {
	lock := sync.RWMutex{}

	for i := 0; i < 3; i++ {
		// 写锁 是互斥的 第二个循环获取不了锁
		lock.Lock()

		defer lock.Unlock()
		fmt.Println("wlock ", i)

	}
}
