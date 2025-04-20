package main

import (
	"fmt"
	"time"
)

func printNumbers() {
	for i := 0; i <= 5; i++ {
		fmt.Println("the printed numbers are : ", i)
		time.Sleep(100 * time.Millisecond) // Simulate work
	}
}
func main() {
	go printNumbers()     // G1
	go printNumbers()     // G2
	time.Sleep(1 * time.Second)
}
