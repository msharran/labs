package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func heavyJob() int {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(80) // n will be between 0 and 10
	time.Sleep(time.Duration(n) * time.Millisecond)
	return n
}

func worker(id int, job <-chan string, results chan<- string) {
	fmt.Println("Worker", id, "started")
	for j := range job {
		fmt.Println("Worker", id, "processing job", j)
		n := heavyJob()
		results <- "job " + j + " done in " + strconv.Itoa(n) + "ms"
	}
}

func main() {
	job := make(chan string)
	results := make(chan string)
	defer close(job)
	defer close(results)

	for i := 0; i < 3; i++ {
		go func(id int) {
			worker(id, job, results)
		}(i)
	}

	go func() {
		for {
			select {
			case results := <-results:
				fmt.Println(results)
			case <-time.After(time.Second * 10):
				fmt.Println("Exiting after 10 seconds")
				os.Exit(1)
			}
		}
	}()

	for i := 1; i <= 2; i++ {
		job <- "job" + strconv.Itoa(i)
	}
}
