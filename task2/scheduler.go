package main

import (
	"fmt"
	"time"
)

func task1() {
	start := time.Now()
	for i := 1; i <= 10; i++ {
		if i%2 != 0 {
			fmt.Println("odd routine: ", i)
			time.Sleep(3 * time.Millisecond)
		}
	}

	fmt.Println("task1 excute time: ", time.Since(start))
}

func task2() {
	start := time.Now()
	for i := 1; i <= 10; i++ {
		if i%2 == 0 {
			fmt.Println("odd routine: ", i)
			time.Sleep(2 * time.Millisecond)
		}
	}

	fmt.Println("task2 excute time: ", time.Since(start))
}

func main() {
	start := time.Now()
	go task1()
	go task2()
	fmt.Println("MainThread .....")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("MainThread excute:", time.Since(start))

}
