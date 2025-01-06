package main

import (
	"fmt"
	"time"
)

func boring(total int) chan string {
	result := make(chan string)
	for i := 0; i < total; i++ {
		go func(i int) {
			time.Sleep(10 * time.Millisecond)
			result <- fmt.Sprintf("boring work %d done", i+1)
		}(i)
	}
	return result
}

func main() {
	fmt.Println("starting my boring task")

	start := time.Now()
	defer func() { fmt.Printf("elapsed: %s", time.Since(start)) }()
	result := boring(10)

	timeout := time.After(10 * time.Millisecond)
	for i := 0; i < 10; i++ {
		select {
		case res := <-result:
			fmt.Println(res)
		case <-timeout:
			fmt.Println("timeout")
			return
		}
	}
	fmt.Println("all done")
}
