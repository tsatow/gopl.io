package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

type client struct {
	Name string
	Ch chan<- string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli.Ch <- msg
			}
		case cli := <-entering:
			if len(clients) > 0 {
				cli.Ch <- "現在チャットには以下のメンバーが参加しています。"
				for member := range clients {
					cli.Ch <- "Name: " + member.Name
				}
			} else {
				cli.Ch <- "あなたがこのチャットの最初のメンバーです。"
			}
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli.Ch)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	// これでいいのかな...?
	input := bufio.NewScanner(conn)
	fmt.Fprint(conn, "まず名を名乗れ！: ")
	input.Scan()

	cli := client{input.Text(), ch}

	timeout := 1 * time.Minute
	timer := time.NewTimer(timeout)
	go func() {
		select {
		case <- timer.C:
			conn.Close()
		}
	}()

	go clientWriter(conn, ch)
	ch <- "You are " + cli.Name
	messages <- cli.Name + " has arrived."
	entering <- cli

	for input.Scan() {
		timer.Reset(timeout)
		messages <- cli.Name + ": " + input.Text()
	}

	leaving <- cli
	messages <- cli.Name + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
