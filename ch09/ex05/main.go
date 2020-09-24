package main

import (
	"fmt"
	"time"
)

func main() {
	ping := make(chan struct{})
	pong := make(chan struct{})
	var count int

	after := time.After(1*time.Second)
	go func() {
		// starter
		ping <- struct{}{}
	}()
	go func(){
		for {
			<- ping
			count++
			pong <- struct{}{}
		}
	}()
	go func() {
		for {
			<-pong
			count++
			ping <- struct{}{}
		}
	}()

	select {
	case <-after:
		fmt.Printf("%d pingpongs\n", count)
	}
}