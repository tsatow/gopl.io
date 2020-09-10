package main

import (
	"../../links"
	"fmt"
	"log"
	"os"
	"sync"
)

const maxdepth = 3


type worklist struct {
	Depth int
	List []string
}

type unseenlink struct {
	Depth int
	Url   string
}

func main() {
	worklists := make(chan worklist)
	unseenLinks := make(chan unseenlink)
	// 処理すべきworklistsとunseenLinksの数をカウントする。
	// worklistsとunseenLinksを個別にみても終了できるタイミングは判断できないと思われるので合計数で管理する。
	var wg sync.WaitGroup

	// 追加するworklistsの分
	wg.Add(1)
	go func() { worklists <- worklist{
		Depth: 0,
		List: os.Args[1:],
	}}()

	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				fmt.Println(link.Url)
				if link.Depth < maxdepth {
					foundLinks := crawl(link.Url)
					// 追加するworklistsの分
					wg.Add(1)
					go func() {
						worklists <- worklist{link.Depth + 1,foundLinks}
					}()
				}
				// 処理したunseenLinksの分
				wg.Done()
			}
		}()
	}

	seen := make(map[string]bool)
	go func() {
		for list := range worklists {
			for _, link := range list.List {
				if !seen[link] {
					// 追加するunseenLinksの分
					wg.Add(1)

					seen[link] = true
					unseenLinks <- unseenlink{list.Depth, link}
				}
			}
			// 処理したworklistsの分
			wg.Done()
		}
	}()

	wg.Wait()
}


func crawl(url string) []string {
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}
