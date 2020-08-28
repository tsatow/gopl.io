package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	echo2()
	fmt.Printf("duration: %v", time.Since(start).Nanoseconds())

	start = time.Now()
	echo3()
	fmt.Printf("duration: %v", time.Since(start).Nanoseconds())
}

func echo2() {
	s, sep := "", ""
	for _, arg := range os.Args[1:] {
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
}

func echo3() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}
