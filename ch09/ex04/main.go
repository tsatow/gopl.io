package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	pipes := 0
	cancel := make(chan struct{})
	tick := time.Tick(1*time.Second)
	go func() {
		var in, out chan struct{}
		in =  make(chan struct{})
		go func() {
			in <- struct{}{}
		}()
		for {
			out = make(chan struct{})
			go func(i <-chan struct{}, o chan<- struct{}) {
				<- i
				pipes++
				o <- struct{}{}
			}(in, out)
			in = out
		}
	}()

	go func() {
		os.Stdin.Read(make([]byte, 1))
		close(cancel)
	}()

	var sec int
	for {
		select {
		case <-tick:
			sec++
			_, _ = fmt.Fprintf(os.Stdout, "%d sec, pipeline size: %d\r", sec, pipes)
		case <- cancel:
			fmt.Println("cancelled")
			return
		}
	}
}