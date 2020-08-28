package main

import (
	"fmt"
	"os"
	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
		os.Exit(1)
	}
	for tag, count := range visit(make(map[string]int), doc) {
		fmt.Printf("%s: %d\n", tag, count)
	}
}

func visit(counts map[string] int, n *html.Node) map[string] int {
	if n.Type == html.ElementNode {
		counts[n.Data] += 1
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		counts = visit(counts, c)
	}
	return counts
}