package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
)

func main() {
	// charcountから副作用を追い出したかったが、
	// 時間もないので妥協しました
	charCount(os.Stdin, os.Stdout)
}

func charCount(reader io.Reader, writer io.Writer) {
	counts := make(map[rune]int)
	var utflen [utf8.UTFMax]int
	invalid := 0

	in := bufio.NewReader(reader)

	for {
		r, n, err := in.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(writer, "charcount: %v\n", err)
			os.Exit(1)
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		counts[r]++
		utflen[n]++
	}
	fmt.Fprintf(writer, "rune\tcount\n")
	for c, n := range counts {
		fmt.Fprintf(writer, "%q\t%d\n", c, n)
	}
	fmt.Fprintf(writer, "\nlen\tcount\n")
	for i, n := range utflen {
		if i > 0 {
			fmt.Fprintf(writer, "%d\t%d\n", i, n)
		}
	}
	if invalid > 0 {
		fmt.Fprintf(writer, "\n%d invalid UTF-8 characters\n", invalid)
	}
}