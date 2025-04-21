package main

import "fmt"

// As a method signature
type Speeker interface {
	bio()
}

type Dog struct {
	breed string
	color string
	name  string
}

func (d *Dog) bio() {
	fmt.Printf("This is %s from %s dreed , it has %s color", d.name, d.breed, d.color)
}

type Human struct {
	name string
}

func (h *Human) bio() {
	fmt.Printf("This is %s .", h.name)
}

func Biography(s Speeker) {
	s.bio()
}

func main() {
	hum := Human{
		name: "jack",
	}
	Biography(&hum)

	dog := Dog{
		breed: "nobel",
		name:  "tommy",
		color: "black",
	}

	Biography(&dog)
}
