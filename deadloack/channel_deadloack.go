package main

import (
	"fmt"
	"time"
)

func deadLockGoroutine() {
	ch := make(chan int, 1) // unbuffer channel
	go func() {
		ch <- 42 // (A) blocks waiting for a receive
		fmt.Println("send data!")
	}()

	//time.Sleep(1 * time.Second)
	select {} // (B) blocks forever, but doesn’t receive
}

// SO WHAT HAPPEND HERE?
// 1. at (B) select doesn't have any cases so it blocked infinitely . Main goroutine never leaves that line.
// 2. meanwhile the child goroutine (A) blocked trying to send on ch because there is no receiver.
// 3. Now all goroutines are blocked:

// a. main is stuck in the empty select {}

// b. sender is stuck on ch <- 42

// 4. The Go runtime notices “hey, nobody can make progress,” and panics

// How to fix it
// You need to pair every send on an unbuffered channel with a matching receive. For example:

func fixDeadlockForUnbufferCh() {
	ch := make(chan int)

	go func() {
		ch <- 42
		fmt.Println("ch receive data!")
	}()

	var result int
	result = <-ch // result receive from ch

	fmt.Println(result)
}

// if you really want to keep main alive without actively receiving, you can use a buffered channel:

func fixDeadlockForBufferCh() {
	ch := make(chan int, 1) // buffer size 1

	go func() {
		ch <- 42 // will succeed immediately
		time.Sleep(1 * time.Second)
		fmt.Println("send data!")
	}()

	select {} // blocks main forever, but no deadlock
}

