package main

import (
	"fmt"
	"sync"
)

func main() {
	var once sync.Once
	s := NewSliceNum()
	
	// why, look once source code: once.done 原子标记是否已处理过
	once.Do(func() {
		s.Add(16)
	})
	once.Do(func() {
		s.Add(16)
	})
	once.Do(func() {
		s.Add(16)
	})
}

type (
	SliceNum []int
)

func (s *SliceNum) Add(elem int) *SliceNum {
	*s = append(*s, elem)
	fmt.Println("add", elem)
	fmt.Println("end adding SliceNum", s)
	return s
}

func NewSliceNum() SliceNum {
	return make(SliceNum, 0)
}