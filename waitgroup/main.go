package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup) {
	fmt.Printf("worker starts %d id \n", id)
	time.Sleep(1 * time.Second)
	fmt.Printf("worker ends with id %d \n", id)
	defer wg.Done()
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i <= 5; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}
	fmt.Println("running routines.....")
	wg.Wait()
	fmt.Println("Done running the goroutines")
}
