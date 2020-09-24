package main

import (
	"fmt"
	"os"
	"sync"
)

func main() {
	_, _ = fmt.Fprintf(os.Stdout, "%d: %d\n", 10, PopCount(10))
	_, _ = fmt.Fprintf(os.Stdout, "%d: %d\n", 20, PopCount(20))
	_, _ = fmt.Fprintf(os.Stdout, "%d: %d\n", 32, PopCount(31))
}

var mkTableOnce sync.Once
var pc [256]byte

func mkTable() {
	fmt.Println("initializing table...")
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
	fmt.Println("initializing table is done.")
}

func PopCount(x uint64) int {
	mkTableOnce.Do(mkTable)
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}