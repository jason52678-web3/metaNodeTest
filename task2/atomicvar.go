package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func count(rid int, cnt *atomic.Int32) {
	for i := 0; i < 1000; i++ {
		cnt.Add(1)
		fmt.Println("routine id:", rid, "\tcounter:", cnt.Load())
	}

}
func main() {

	cnt := atomic.Int32{}

	for i := 0; i < 10; i++ {
		go count(i+1, &cnt)
	}

	time.Sleep(1 * time.Second)
	fmt.Println(cnt.Load())

}
