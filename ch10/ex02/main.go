package main

import (
	"./archive"
	_ "./archive/tar"
	_ "./archive/zip"
	"flag"
	"fmt"
	"os"
)

var f = flag.String("f", "", "file")

func main() {
	if *f == "" {
		fmt.Fprint(os.Stderr, "fileを指定してください。")
		os.Exit(1)
	}

	file, err := os.Open(*f)
	if err != nil {
		fmt.Fprint(os.Stderr, "fileの読み込みに失敗しました。")
		os.Exit(1)
	}

}