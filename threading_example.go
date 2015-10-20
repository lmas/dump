package main

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

func main() {
	workers := 3
	max_jobs := 999999
	var i int
	var group sync.WaitGroup
	signaller := make(chan bool)
	queue := make(chan string)

	// init the workers
	fmt.Println("Starting workers...")
	for i = 1; i <= workers; i++ {
		group.Add(1)
		go worker(i, &group, signaller, queue)
	}

	// add some work to the queue
	fmt.Println("Adding jobs to queue...")
	for i = 1; i <= max_jobs; i++ {
		queue <- fmt.Sprintf("msg #%d", i)
	}

	// wait for workers to finish their job
	fmt.Println("Waiting for workers to finish...")
	close(signaller)
	group.Wait()

	// done
	fmt.Println("Exit.")
}

func worker(id int, group *sync.WaitGroup, signaller chan bool, queue chan string) {
	defer group.Done()
	jobs := 0
	fmt.Printf("#%d: started\n", id)

	for {
		select {
		case <-signaller:
			fmt.Printf("#%d: signal, closing. Did %d jobs.\n", id, jobs)
			return
			//case <-queue:
		case msg := <-queue:
			tmp := []byte(msg)
			sha256.Sum256(tmp)
			//fmt.Printf("#%d: %s (%s)\n", id, msg, hash)
			jobs += 1
		}
	}
}
