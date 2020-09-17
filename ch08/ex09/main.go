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

	errors := make(chan error)
	var errs []error

	initializeScreen(roots)

	var wg sync.WaitGroup
	for i, root := range roots {
		wg.Add(1)
		go func(idx int, rootDir string) {
			defer wg.Done()

			fileSizes := make(chan int64)
			var wgByRoot sync.WaitGroup

			wgByRoot.Add(1)
			go walkDir(rootDir, &wgByRoot, fileSizes, errors)
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

	go func() {
		for err := range errors {
			errs = append(errs, err)
		}
	}()

	wg.Wait()
	_, _ = fmt.Fprintf(os.Stdout, "\033[%dB", len(roots))

	for _, err := range errs {
		_, _ = fmt.Fprintf(os.Stderr, "\u001b[00;31merror: %v\u001b[00m\n", err)
	}
}

func initializeScreen(roots []string) {
	for range roots {
		_, _ = fmt.Fprint(os.Stdout, "\n")
	}
	_, _ = fmt.Fprintf(os.Stdout, "\033[%dA", len(roots))
}

func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64, errors chan<- error) {
	defer n.Done()
	for _, entry := range dirents(dir, errors) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, n, fileSizes, errors)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

var sema = make(chan struct{}, 20)

func dirents(dir string, errors chan<- error) []os.FileInfo {
	sema <- struct{}{}
	defer func() { <-sema }()

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		errors <- err
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
