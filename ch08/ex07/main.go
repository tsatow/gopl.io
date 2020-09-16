package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type content struct {
	Path  string
	Bytes []byte
}

type context struct {
	Site     string
	Hostname string
	RootDir  string
}

var urlReplacer = strings.NewReplacer("http://", "./", "https://", "./")

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "一個のサイトを指定してください。\n")
	}

	var wg sync.WaitGroup
	site := strings.TrimSpace(os.Args[1])
	parsedUrl, err := url.Parse(site)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "解析できないURL: %s\n", site)
		os.Exit(1)
	}
	if !parsedUrl.IsAbs() {
		_, _ = fmt.Fprintf(os.Stderr, "絶対URLで指定してください。: %s\n", site)
		os.Exit(1)
	}
	// MkdirAllだとディレクトリ名に`.`が含まれるとエラーになるのでアドホックに対応
	rootPath := filepath.Join(".", parsedUrl.Hostname())
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		if err := os.Mkdir(filepath.Join(".", parsedUrl.Hostname()), 0777); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "rootフォルダの作成に失敗しました。: %s\n", parsedUrl.Hostname())
			os.Exit(1)
		}
	}
	context := context{site, parsedUrl.Hostname(), rootPath}

	preResourceGetRequests := make(chan []string)
	resourceGetRequests := make(chan string)
	contentWriteRequests := make(chan content)

	wg.Add(1)
	go func() {
		preResourceGetRequests <- []string{site}
	}()

	seen := make(map[string]bool)
	go func() {
		for links := range preResourceGetRequests {
			for _, link := range links {

				if !seen[link] {
					wg.Add(1)

					seen[link] = true
					resourceGetRequests <- link
				}
			}
			wg.Done()
		}
	}()

	for i := 0; i < 10; i++ {
		go func() {
			for resource := range resourceGetRequests {
				getResource(context, contentWriteRequests, preResourceGetRequests, resource, &wg)
				wg.Done()
			}
		}()
	}

	go writeFile(context, contentWriteRequests, &wg)

	go func(delay time.Duration) {
		for {
			for _, r := range `-\|/` {
				fmt.Printf("\r%c Scanning...\r", r)
				time.Sleep(delay)
			}
		}
	}(100 * time.Millisecond)

	wg.Wait()
}

func getResource(ctx context, contentWriteRequests chan<- content, preResourceGetRequests chan<- []string, url string, wg *sync.WaitGroup) {
	resp, err := http.Get(url)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "リソースの取得に失敗しました。: %s\n", url)
		return
	}
	if resp == nil {
		// 意味ないと思うけど警告避け
		_, _ = fmt.Fprintf(os.Stderr, "レスポンスが空でした。: %s\n", url)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_, _ = fmt.Fprintf(os.Stderr, "異常なステータス %s: %s\n", url, resp.Status)
	}

	if isHtml(&resp.Header) {
		// ほとんどエラーにならないのでここでは無視する
		node, _ := html.Parse(resp.Body)
		links := visit(url, ctx, nil, node)

		wg.Add(1)
		go func() {
			preResourceGetRequests <- links
		}()
		wg.Add(1)
		go func() {
			contentWriteRequests <- content{completeIndexHtml(ctx, absUrlToPath(url)), toBytes(node)}
		}()
	} else {
		body, err := readAllAsBytes(resp.Body)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "レスポンスボディの読み込みに失敗しました。: %s\n", url)
			return
		}

		wg.Add(1)
		go func() {
			contentWriteRequests <- content{absUrlToPath(url), body}
		}()
	}
}

func visit(currUrl string, ctx context, links []string, n *html.Node) []string {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "a", "link":
			for i, attr := range n.Attr {
				if attr.Key == "href" {
					completedUrl := completeUrl(ctx, currUrl, attr.Val)

					if strings.TrimSpace(completedUrl.Hostname()) == strings.TrimSpace(ctx.Hostname) {
						completedUrlStr := completedUrl.String()
						links = append(links, removeFragmentAndQuery(*completedUrl).String())

						// completeIndexHtmlはHTMLと確定していないとうまく動かないので、、、
						if strings.HasSuffix(completedUrlStr, "/") {
							completedUrlStr += "index.html"
						}
						if relativePath, err := absUrlToRelativePath(currUrl, completedUrlStr); err == nil {
							n.Attr[i] = html.Attribute{Namespace: attr.Namespace, Key: attr.Key, Val: relativePath}
						}
					}
				}
			}
		case "img", "script":
			for i, attr := range n.Attr {
				if attr.Key == "src" {
					completedUrl := completeUrl(ctx, currUrl, attr.Val)

					if strings.TrimSpace(completedUrl.Hostname()) == strings.TrimSpace(ctx.Hostname) {
						completedUrlStr := completedUrl.String()
						links = append(links, removeFragmentAndQuery(*completedUrl).String())

						// completeIndexHtmlはHTMLと確定していないとうまく動かないので、、、
						if strings.HasSuffix(completedUrlStr, "/") {
							completedUrlStr += "index.html"
						}
						if relativePath, err := absUrlToRelativePath(currUrl, completedUrlStr); err == nil {
							n.Attr[i] = html.Attribute{Namespace: attr.Namespace, Key: attr.Key, Val: relativePath}
						}
					}
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(currUrl, ctx, links, c)
	}

	return links
}

func writeFile(ctx context, contentWriteRequest <-chan content, wg *sync.WaitGroup) {
	// 競合なく並行で処理させるのが大変そうなので単一のゴルーチンで処理する
	for content := range contentWriteRequest {
		dir := filepath.Dir(content.Path)
		if err := os.MkdirAll(dir, 0777); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "ディレクトリの作成に失敗しました。: dir:%s, contentpath:%s 原因: %v\n", dir, content.Path, err)
		} else {
			// rootDirは最初に作成済
			if content.Path != ctx.RootDir {
				if file, err := os.Create(content.Path); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "ファイルの作成に失敗しました。: %s 原因: %v\n", dir, err)
				} else {
					if err := file.Chmod(0777); err != nil {
						_, _ = fmt.Fprintf(os.Stderr, "ファイルの権限変更に失敗しました。: %s 原因: %v\n", dir, err)
					}
					if _, err := file.Write(content.Bytes); err != nil {
						_, _ = fmt.Fprintf(os.Stderr, "ファイルの書き出しに失敗しました。: %s 原因: %v\n", dir, err)
					}
					_ = file.Close()
				}
			}
		}
		wg.Done()
	}
}

func isHtml(respHeader *http.Header) bool {
	return strings.Contains(respHeader.Get("Content-Type"), "text/html")
}

func toBytes(node *html.Node) []byte {
	buf := bytes.NewBufferString("")
	html.Render(buf, node)
	return buf.Bytes()
}

func readAllAsBytes(in io.ReadCloser) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, in)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func completeUrl(ctx context, currUrl, rawUrl string) *url.URL {
	// ロジック全体的にだめだめな気がする

	u, err := url.Parse(rawUrl)
	if err != nil {
		// この場合はurl packageでは解析できない
		if strings.HasPrefix(rawUrl, "/") {
			// 一度解析されているはずなのでエラーにならない
			baseUrl, _ := url.Parse(ctx.Site)
			baseUrl.Path = path.Join(baseUrl.Path, rawUrl[1:])
			if path.Dir(baseUrl.Path) == baseUrl.Path {
				baseUrl.Path = path.Join(baseUrl.Path, "index.html")
			}
			return baseUrl
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "解析できないURL: %s\n", rawUrl)
			return nil
		}
	}
	// 最初にパースしてるのでここではエラーが発生しない想定
	baseUrl, _ := url.Parse(currUrl)
	completeUrl := baseUrl.ResolveReference(u)

	if path.Dir(completeUrl.Path) == completeUrl.Path {
		completeUrl.Path = path.Join(completeUrl.Path, "index.html")
	}
	return completeUrl
}

func completeIndexHtml(ctx context, p string) string {
	if path.Dir(p) == p {
		return path.Join(p, "index.html")
	}
	if ctx.RootDir == p {
		return path.Join(p, "index.html")
	}
	if !strings.HasSuffix(p, ".html") && !strings.HasSuffix(p, ".htm") {
		return path.Join(p, "index.html")
	}
	return p
}

func removeFragmentAndQuery(u url.URL) *url.URL {
	u.ForceQuery = false
	u.RawQuery = ""
	u.Fragment = ""

	return &u
}

func absUrlToPath(absUrl string) string {
	return filepath.Join(".", urlReplacer.Replace(absUrl))
}

func absUrlToRelativePath(srcAbsUrl, destAbsUrl string) (string, error) {
	relativePath, err := filepath.Rel(absUrlToPath(srcAbsUrl), absUrlToPath(destAbsUrl))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "相対パスの計算に失敗しました: %s, %s\n", absUrlToPath(srcAbsUrl), absUrlToPath(destAbsUrl))
		// 算出不可能なときはURLのままにしちゃう
		return "", err
	}

	return relativePath, nil
}
