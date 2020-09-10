package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Fprint(os.Stdout, "wait...\n\n\033[2A")

	for index, arg := range os.Args[1:] {
		arr := strings.Split(arg, "=")

		if len(arr) != 2 {
			fmt.Fprint(os.Stderr, "invalid argument.")
			os.Exit(1)
		}

		city := arr[0]
		url := arr[1]

		go dialTime(index, city, url)
	}

	select {}
}

func dialTime(index int, city string, url string) {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	cityClock(index, city, conn)
}

func cityClock(index int, city string, src io.Reader) {
	scanner := bufio.NewScanner(src)

	for scanner.Scan() {
		time := strings.TrimSpace(scanner.Text())
		if index == 0 {
			fmt.Fprintf(os.Stdout, "\r%s: %s\r", city, time)
		} else {
			fmt.Fprintf(os.Stdout, "\r\033[%dB%s: %s\033[%dA\r", index, city, time, index)
		}
	}
}
