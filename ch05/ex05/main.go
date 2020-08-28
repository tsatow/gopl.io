package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("give me an url!!!")
	}

	url := os.Args[1]
	words, images, err := CountAndImages(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}

	fmt.Printf("site %s has %d words and %d images.\n", url, words, images)
}

func CountAndImages(url string) (words, images int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	doc, err := html.Parse(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("parsing HTML: %s", err)
		return
	}
	words, images = countAndImages(doc)
	return
}

func countAndImages(n *html.Node) (words, images int) {
	if n.Type == html.ElementNode && n.Data == "img" {
		images++
	} else if n.Type == html.TextNode {
		scanner := bufio.NewScanner(strings.NewReader(n.Data))
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			words++
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		childNodeWords, childNodeImages := countAndImages(c)
		words += childNodeWords
		images += childNodeImages
	}

	return
}