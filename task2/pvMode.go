package main

import (
	"fmt"
	"math/rand"
	"time"
)

func producer(ch chan<- int) {
	for {
		data := rand.Intn(999)
		ch <- data
		fmt.Println("producer data: ", data)
		time.Sleep(10 * time.Millisecond)
		if len(ch) == cap(ch) {
			fmt.Println("channel已经满，等待消费者消费")
		}
	}
}

func consumer(ch <-chan int) {

	for value := range ch {
		fmt.Println("Comsumer get :", value)
		time.Sleep(200 * time.Millisecond)
	}
}

func main() {
	//ch := make(chan int, 5)
	ch := make(chan int, 100)

	go producer(ch)
	go consumer(ch)

	for {
		select {
		case _, ok := <-ch:
			if !ok {
				fmt.Println("channel closed")
				return
			}
		case <-time.After(1 * time.Second):
			fmt.Println("timeout")
		default:
			fmt.Println(cap(ch), "Channel为空，等待生产者生产。。。")
			time.Sleep(1 * time.Second)
		}
	}

}
