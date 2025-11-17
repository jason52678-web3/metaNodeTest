package main

import (
	"fmt"
	"time"
)

func getOldNumber() {
	for i := 1; i <= 10; i++ {
		if i%2 != 0 {
			fmt.Println("odd routine: ", i)
			time.Sleep(30 * time.Millisecond)
		}
	}
}

func getEvenNumber() {
	for i := 1; i <= 10; i++ {
		if i%2 == 0 {
			fmt.Println("even routine: ", i)
		}
		time.Sleep(50 * time.Millisecond)
	}
}
func main() {
	go getOldNumber()
	go getEvenNumber()

	for i := 0; i < 10; i++ {
		fmt.Println("MainThread: ", i)
		time.Sleep(100 * time.Millisecond)
	}
}
