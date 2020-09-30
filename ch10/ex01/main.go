package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
)

var format = flag.String("format", "jpeg", "output format")

func main() {
	img, kind, err := image.Decode(os.Stdin)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "Input format =", kind)
	if *format == "jpeg" {
		jpeg.Encode(os.Stdout, img, &jpeg.Options{Quality: 95})
	} else if *format == "gif" {
		gif.Encode(os.Stdout, img, &gif.Options{})
	} else if *format == "png" {
		png.Encode(os.Stdout, img)
	} else {
		fmt.Fprintf(os.Stderr, "unsupported format %s", *format)
		os.Exit(1)
	}
}