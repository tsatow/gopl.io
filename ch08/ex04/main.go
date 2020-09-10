package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		var wg sync.WaitGroup
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		wg.Add(1)
		go func (c net.Conn) {
			defer wg.Done()
			input := bufio.NewScanner(c)
			for input.Scan() {
				echo(c, input.Text(), 1*time.Second)
			}
		}(conn)

		wg.Wait()

		if tcpConn, ok := conn.(*net.TCPConn); ok {
			tcpConn.CloseWrite()
		}
	}
}

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}