package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

type client struct {
	c    chan<- string //对外发送消息的通道
	name string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // 所有接受的客户消息
)

const timeout = 5 * time.Second

func main() {
	listenner, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
		return
	}
	go boardCaster()
	for {
		conn, err := listenner.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	input := bufio.NewScanner(conn)
	who := conn.RemoteAddr().String()
	msg := make(chan string)
	client := client{msg, who}
	go func() {
		for msg := range msg {
			fmt.Fprint(conn, msg)
		}
	}()
	entering <- client
	fmt.Fprintf(conn, "you are %s\n", conn.RemoteAddr().String())

	send := make(chan string)
	go func() {
		for {
			select {
			case <-time.After(timeout):
				conn.Close()
			case m := <-send:
				messages <- m
			}
		}
	}()
	for input.Scan() {
		send <- fmt.Sprintf("%s: %s\n", who, input.Text())
	}
	defer func() {
		conn.Close()
		leaving <- client
	}()
}

func boardCaster() {
	clients := make(map[client]bool) //所有连接的客户端
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli.c <- msg
			}
		case cli := <-entering:
			go func() {
				messages <- fmt.Sprintf("%s: enter\n", cli.name)
			}()
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			go func() {
				messages <- fmt.Sprintf("%s: leave\n", cli.name)
			}()
			close(cli.c)
		}
	}
}
