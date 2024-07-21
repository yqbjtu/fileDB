package main

import (
	"fmt"
	"sync"
)

var pool *sync.Pool

type Person struct {
	Name string
}

func initPool() {
	pool = &sync.Pool{
		New: func() interface{} {
			return new(Person)
		},
	}
}

func main() {
	initPool()

	p := pool.Get().(*Person)
	p.Name = "first"

	// 需要在 Put 前，清除 p 的各个成员，这里为了展示，就不清除了
	// 放回对象池中以供其他goroutine调用
	pool.Put(p)
	fmt.Println("Get from pool:", pool.Get().(*Person))
	p2 := pool.Get().(*Person)
	fmt.Println("Pool is empty", p2)
	if nil == p2 {
		fmt.Println("p2 is nil")
	}
}
