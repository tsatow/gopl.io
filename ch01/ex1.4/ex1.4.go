package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]Counts)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts)
			f.Close()
		}
	}

	for line, c := range counts {
		if c.counts > 1 {
			fmt.Printf("%d\t%s\n", c.counts, line)
			for _, fileName := range c.files {
				fmt.Println(fileName)
			}
		}
	}
}

func countLines(f *os.File, counts map[string]Counts) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		var c = counts[input.Text()]
		c.counts = c.counts + 1
		c.files = append(c.files, f.Name())

		counts[input.Text()] = c
	}
}

type Counts struct {
	files  []string
	counts int
}
