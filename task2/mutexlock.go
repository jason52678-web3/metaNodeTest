package main

import (
	"fmt"
	"sync"
	"time"
)

var lock sync.Mutex

func counter(rid int, cnt *int) {
	for i := 0; i < 1000; i++ {
		lock.Lock()
		*cnt++
		fmt.Println("routine id:", rid, "\tcounter:", *cnt)
		lock.Unlock()
	}
}
func main() {
	cnt := 0

	for i := 0; i < 10; i++ {
		go counter(i+1, &cnt)
	}

	time.Sleep(1 * time.Second)
	fmt.Println(cnt)
}
