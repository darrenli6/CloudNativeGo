package main

import (
	"fmt"
	"sync"
)

// 单利模式

type SliceNum []int

func NewSlice() SliceNum {
	return make(SliceNum, 0)
}

func (s *SliceNum) add(elem int) *SliceNum {

	*s = append(*s, elem)
	fmt.Println("add", elem)
	fmt.Println("add sliceNum end", s)
	return s
}

func main() {
	var once sync.Once

	s := NewSlice()

	// 只会调用一次

	once.Do(func() {
		s.add(13)
	})

	once.Do(func() {
		s.add(13)
	})

	once.Do(func() {
		s.add(13)
	})

}
