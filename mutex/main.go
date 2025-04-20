package main

import (
	"fmt"
	"sync"
)

type SafeCounter struct {
	mutex sync.Mutex
	count int
}

func (sc *SafeCounter) counter() {
	sc.mutex.Lock() // active lock
	sc.count++
	fmt.Println(sc.count)
	defer sc.mutex.Unlock() // release lock
}

func main() {
	counter := SafeCounter{}

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			counter.counter()
			defer wg.Done()
		}()

	}

	wg.Wait()
	fmt.Println("done!!")
}
