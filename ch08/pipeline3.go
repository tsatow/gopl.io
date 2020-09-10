package main

import (
	"fmt"
	"time"
)

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	go counter(naturals)
	go square(naturals, squares)
	printer(squares)
}

func counter(out chan<- int) {
	for x := 0; x < 100; x++ {
		time.Sleep(100 * time.Millisecond)
		out <- x
	}
	close(out)
}

func square(in <-chan int, out chan<- int) {
	for x := range in {
		out <- x * x
	}
	close(out)
}

func printer(in <-chan int) {
	for x := range in {
		fmt.Println(x)
	}
}
