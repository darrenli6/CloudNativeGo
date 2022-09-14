package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	go rlock()
	go wlock()
	go lock()

	time.Sleep(5 * time.Second)
}

func lock() {
	lock := sync.Mutex{}

	for i := 0; i < 3; i++ {
		lock.Lock()
		defer lock.Unlock()
		fmt.Println("lock ", i)

	}
}

func rlock() {
	lock := sync.RWMutex{}

	for i := 0; i < 3; i++ {
		lock.RLock()
		defer lock.RUnlock()
		fmt.Println("rlock ", i)

	}
}

func wlock() {
	lock := sync.RWMutex{}

	for i := 0; i < 3; i++ {
		lock.Lock()
		defer lock.Unlock()
		fmt.Println("wlock ", i)

	}
}
