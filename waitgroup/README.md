# sync.WaitGroup
A WaitGroup is used to wait for a collection of goroutines to finish executing.

## How it works:
    1. Add(): Sets the number of goroutines to wait for

    2. Done(): Decrements the counter (called when a goroutine finishes)

    3. Wait(): Blocks until the counter reaches zero


## Key points:
    1. Always pass WaitGroup by pointer (it contains state)

    2. Call Add() before starting the goroutine

    3. Use defer wg.Done() to ensure it's called even if the goroutine panics

    4. Wait() blocks until all Done() calls complete

