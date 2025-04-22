# Deadlock
### Definition :
A deadlock arises when two or more goroutines are each waiting for the other to release a resource (e.g., a lock or channel), so none can proceed. In Go, deadlocks often involve unbuffered channels or mutexes.

## Classic Symptoms in Go
1. Program hangs indefinitely.
2. The Go runtime may panic:
     fatal error: all goroutines are asleep - deadlock!

## Example 1: Channel-based Deadlock
```go
package main

func main() {
    ch := make(chan int) // unbuffered

    // Sender goroutine
    go func() {
        ch <- 42          // blocks, since no receiver yet
        fmt.Println("sent")
    }()

    // Receiver goroutine
    go func() {
        fmt.Println("received:", <-ch) // blocks, since send hasn't happened
    }()

    select {} // block forever
}
```

## Example 2: Mutex-based Deadlock
### What is a Mutex?
 A mutex (mutual exclusion) is used to protect shared resources from concurrent access.
### What is a Mutex-based Deadlock?
A mutex-based deadlock happens when:

1. Two (or more) goroutines are each holding a lock,and each is waiting for the other goroutine to release another lock.

2. Since both are waiting, neither ever proceeds = deadlock.


```go
package main

import (
    "sync"
)

func main() {
    var mu1, mu2 sync.Mutex

    go func() {
        mu1.Lock()
        defer mu1.Unlock()

        // simulate work
        mu2.Lock()    // tries to acquire mu2
        defer mu2.Unlock()
    }()

    go func() {
        mu2.Lock()
        defer mu2.Unlock()

        // simulate work
        mu1.Lock()    // tries to acquire mu1
        defer mu1.Unlock()
    }()

    select {} // block
}
```


# How to Prevent Deadlock
1. Lock Ordering: Always acquire multiple locks in a consistent order.

2. Channel Buffering / Select with Default:

3. Use buffered channels to decouple sends and receives.

4. Use select with a default case to avoid blocking if nobody is ready.

5. Timeouts / Contexts: Use context.Context with deadlines or timeouts.

6. Avoid Circular Wait: Design so goroutines don’t wait on each other in a cycle.
