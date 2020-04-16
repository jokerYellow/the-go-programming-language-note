package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client chan<- string //对外发送消息的通道

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // 所有接受的客户消息
)

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
	client := make(chan string)
	who := conn.RemoteAddr().String()
	entering <- client
	messages <- fmt.Sprintf("%s: enter\n", who)
	go func() {
		for msg := range client {
			conn.Write([]byte(msg))
		}
	}()
	fmt.Fprintf(conn, "you are %s\n", conn.RemoteAddr().String())
	for input.Scan() {
		messages <- fmt.Sprintf("%s: %s\n", who, input.Text())
	}
	defer func() {
		conn.Close()
		leaving <- client
		messages <- fmt.Sprintf("%s: leave\n", who)
	}()
}

func boardCaster() {
	clients := make(map[client]bool) //所有连接的客户端
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}
