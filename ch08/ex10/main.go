package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
)

func main() {
	worklist := make(chan []string)
	unseenLinks := make(chan string)
	cancel := make(chan struct{})
	// 必要な関数を匿名関数で用意
	cancelled := func() bool {
		select {
		case <-cancel:
			return true
		default:
			return false
		}
	}
	var forEachNode func(n *html.Node, pre, post func(n *html.Node))
	forEachNode = func(n *html.Node, pre, post func(n *html.Node)) {
		if cancelled() {
			return
		}
		if pre != nil {
			pre(n)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if cancelled() {
				return
			}
			forEachNode(c, pre, post)
		}

		if post != nil {
			post(n)
		}
	}
	extract := func(url string) ([]string, error) {
		req, _ := http.NewRequest("GET", url, nil)
		req.Cancel = cancel
		client := new(http.Client)
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
		}

		if cancelled() {
			return nil, errors.New(" Process is cancelled. ")
		}
		doc, err := html.Parse(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
		}

		var links []string
		visitNode := func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "a" {
				for _, a := range n.Attr {
					if cancelled() {
						return
					}
					if a.Key != "href" {
						continue
					}
					link, err := resp.Request.URL.Parse(a.Val)
					if err != nil {
						continue
					}
					links = append(links, link.String())
				}
			}
		}
		forEachNode(doc, visitNode, nil)
		return links, nil
	}
	crawl := func(url string) []string {
		fmt.Println(url)
		list, err := extract(url)
		if err != nil {
			log.Print(err)
		}
		return list
	}

	// ここから処理本体
	go func() { worklist <- os.Args[1:] }()

	go func() {
		os.Stdin.Read(make([]byte, 1))
		close(cancel)
	}()

	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				if cancelled() {
					return
				}
				foundLinks := crawl(link)
				go func() { worklist <- foundLinks }()
			}
		}()
	}

	seen := make(map[string]bool)
	for {
		select {
		case <-cancel:
			close(cancel)
			for range worklist {
			}
			for range unseenLinks {
			}
			return
		case list := <- worklist:
			for _, link := range list {
				if cancelled() {
					return
				}
				if !seen[link] {
					seen[link] = true
				unseenLinks <- link
				}
			}
		}
	}
}
