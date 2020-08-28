package main

import (
	"fmt"
	"os"
	"golang.org/x/net/html"
	"strings"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
		os.Exit(1)
	}
	for _, text := range visit(nil, doc) {
		fmt.Printf("%s\n", text)
	}
}

func visit(texts []string, n *html.Node) []string {
	if n.Type == html.TextNode {
		trimmed := strings.TrimSpace(n.Data)
		if trimmed != "" {
			texts = append(texts, trimmed)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		childTag := strings.ToLower(c.Data)
		if childTag != "script" && childTag != "style" {
			texts = visit(texts, c)
		}
	}
	return texts
}