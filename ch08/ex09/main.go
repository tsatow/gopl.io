package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	initializeScreen(roots)

	var wg sync.WaitGroup
	for i, root := range roots {
		wg.Add(1)
		go func(idx int, rootDir string) {
			defer wg.Done()

			fileSizes := make(chan int64)
			var wgByRoot sync.WaitGroup

			wgByRoot.Add(1)
			go walkDir(rootDir, &wgByRoot, fileSizes)
			go func() {
				wgByRoot.Wait()
				close(fileSizes)
			}()

			tick := time.Tick(500 * time.Millisecond)
			var nfiles, nbytes int64
		loop:
			for {
				select {
				case size, ok := <-fileSizes:
					if !ok {
						break loop
					}
					nfiles++
					nbytes += size
				case <-tick:
					printDiskUsage(idx, rootDir, nfiles, nbytes)
				}
			}
			printDiskUsage(idx, rootDir, nfiles, nbytes)
		}(i, root)
	}

	wg.Wait()
	fmt.Fprintf(os.Stdout, "\033[%dB", len(roots))
}

func initializeScreen(roots []string) {
	for range roots {
		_, _ = fmt.Fprint(os.Stdout, "\n")
	}
	_, _ = fmt.Fprintf(os.Stdout, "\033[%dA", len(roots))
}

func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, n, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

var sema = make(chan struct{}, 20)

func dirents(dir string) []os.FileInfo {
	sema <- struct{}{}
	defer func() { <-sema }()

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ex09: %v", err)
		return nil
	}
	return entries
}

func printDiskUsage(idx int, root string, nfiles, nbytes int64) {
	if idx == 0 {
		fmt.Fprintf(os.Stdout, "\r%s: %d files %.1f GB\r", root, nfiles, float64(nbytes)/1e9)
	} else {
		fmt.Fprintf(os.Stdout, "\r\033[%dB%s: %d files %.1f GB\033[%dA\r", idx, root, nfiles, float64(nbytes)/1e9, idx)
	}
}
