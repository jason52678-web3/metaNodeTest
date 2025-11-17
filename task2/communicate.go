package main

import (
	"fmt"
	"time"
)

func sendMsg(channel chan<- int) {
	for i := 1; i <= 10; i++ {
		channel <- i
		time.Sleep(500 * time.Millisecond)
	}
	close(channel)
}

func recvMsg(channel <-chan int) {
	for value := range channel {
		fmt.Println("recvMsg from channel: ", value)
	}

}
func main() {
	channel := make(chan int)
	go sendMsg(channel)
	go recvMsg(channel)

	for {
		select {
		case _, ok := <-channel:
			if !ok {
				fmt.Println("channel closed")
				return
			}
		default:
			fmt.Println("没有数据，等待中。。。")
			time.Sleep(1 * time.Second)
		}
	}

}
