package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func heavyJob() int {
	// rand.Seed(time.Now().UnixNano())
	// n := rand.Intn(5)
	// time.Sleep(time.Duration(n) * time.Second)
	n := 1
	time.Sleep(time.Duration(n) * time.Second)
	return n
}

func worker(id int, jobs <-chan string) {
	fmt.Println("Worker", id, "started")
	for j := range jobs {
		fmt.Println("Worker", id, "processing job", j)
		n := heavyJob()
		fmt.Println(j + " done in " + strconv.Itoa(n) + "s")
	}
}

func main() {
	start := time.Now()
	job := make(chan string)

	workerswg := &sync.WaitGroup{}
	jobswg := &sync.WaitGroup{}

	for i := 1; i <= 30; i++ {
		workerswg.Add(1)
		go func(i int) {
			defer workerswg.Done()
			worker(i, job)
		}(i)
	}

	jobswg.Add(1)
	go func() {
		defer jobswg.Done()
		for i := 1; i <= 10; i++ {
			job <- "job_a" + strconv.Itoa(i)
		}
	}()

	jobswg.Add(1)
	go func() {
		defer jobswg.Done()
		for i := 1; i <= 10; i++ {
			job <- "job_b" + strconv.Itoa(i)
		}
	}()

	jobswg.Wait()
	close(job) // close channel after sending jobs to workers.
	workerswg.Wait()
	fmt.Println("All jobs done in ", time.Since(start))

}
