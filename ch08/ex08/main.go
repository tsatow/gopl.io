package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		fmt.Println("new connection.")
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func handleConn(c net.Conn) {
	input := bufio.NewScanner(c)
	eof := make(chan struct{})
	timeout := 10 * time.Second
	timer := time.NewTimer(timeout)

	go func() {
		for input.Scan() {
			timer.Reset(timeout)
			echo(c, input.Text(), 1*time.Second)
		}
		eof <- struct{}{}
	}()

	select {
	case <-eof:
		fmt.Println("connection closed.")
		_ = c.Close()
	case <-timer.C:
		fmt.Println("timer is expired.")
		_ = c.Close()
	}

	close(eof)
}
