package main

import (
	"fmt"
	"sync"
	"time"
)

// Circular Waiting

func mutexBasedDeadlock() {
	var mu1, mu2 sync.Mutex

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		mu1.Lock()
		fmt.Println("Goroutine 1: locked mu1")

		// simulate work
		x := 2
		fmt.Println(x)
		time.Sleep(1 * time.Second)

		mu2.Lock()
		fmt.Println("Goroutine 1: locked mu2")

		mu2.Unlock()
		mu1.Unlock()
	}()

	go func() {
		defer wg.Done()
		mu2.Lock()
		fmt.Println("Goroutine 2: locked mu2")

		// simulate work
		y := 2
		fmt.Println(y)
		time.Sleep(1 * time.Second)

		mu1.Lock()
		fmt.Println("Goroutine 2: locked mu1")

		mu1.Unlock()
		mu2.Unlock()
	}()

	wg.Wait()
}

// What’s Happening:

// 1. Goroutine 1 locks mu1, and tries to lock mu2.
// 2. Goroutine 2 locks mu2, and tries to lock mu1.
// Both are now waiting for each other → 💥 Deadlock!
// This is called a classic deadlock pattern.

// Fix: Lock Ordering
// Always acquire multiple locks in the same order in all goroutines.
func fixedMutexBasedDeadlock() {
	var mu1, mu2 sync.Mutex

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		mu1.Lock()
		fmt.Println("Goroutine 1: locked mu1")

		// simulate work
		x := 2
		fmt.Println(x)
		time.Sleep(1 * time.Second)

		mu2.Lock()
		fmt.Println("Goroutine 1: locked mu2")

		mu2.Unlock()
		mu1.Unlock()
	}()

	go func() {
		defer wg.Done()
		mu1.Lock()
		fmt.Println("Goroutine 2: locked mu2")

		// simulate work
		y := 2
		fmt.Println(y)
		time.Sleep(1 * time.Second)

		mu2.Lock()
		fmt.Println("Goroutine 2: locked mu1")

		mu1.Unlock()
		mu2.Unlock()
	}()

	wg.Wait()
}

// Now both goroutines lock in the same order, so one will wait for the first lock and the other will proceed — no deadlock.
