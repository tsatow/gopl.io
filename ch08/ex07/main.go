package main

import (
	"fmt"
	"golang.org/x/net/html"
	"os"
	"path/filepath"
	"sync"
)

type contents struct {
	Path string
	Bytes []byte
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "一個のサイトを指定してください。")
	}

	var wg sync.WaitGroup
	site := os.Args[1]

	contents := make(chan contents)

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

	go writeFile(contents, &wg)
}

func findLinks(url string) []string {
	// 同一ドメインの場合はhttp.Getしてリソースを取得(外部ドメインの場合は何もしない)
	// 1. リソースがHTMLの場合
	// script、img、link、style、aタグからリンクを探す
	// リンクが同一ドメイン内にあるのであれば、相対パスに修正する
	// 修正したHTMLを書き出す(パスの場合はindex.htmlにする)
	// リンクを戻り値として返す、もしくはワークリストとしてchに送信する
	// 2. リソースがHTML以外の場合
	// リソースを適切なパスに書き出す
}

func writeFile(in <-chan contents, wg *sync.WaitGroup) {
	// 競合なく並行で処理させるのが大変そうなので単一のゴルーチンで処理する
	for content := range in {
		dir := filepath.Dir(content.Path)
		if err := os.MkdirAll(dir, 744); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "ディレクトリの作成に失敗しました。: %s", dir)
		} else {

			if file, err := os.Create(content.Path); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "ファイルの作成に失敗しました。: %s", dir)
			} else {
				if _, err := file.Write(content.Bytes); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "ファイルの書き出しに失敗しました。: %s", dir)
				}
				_ = file.Close()
			}
		}
		wg.Done()
	}
}