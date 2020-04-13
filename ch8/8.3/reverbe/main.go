package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func handleConn(conn net.Conn) {
	input := bufio.NewScanner(conn)
	requestTodo := make(chan int)
	closedRead := false
	log.Printf("%s:connected", conn.RemoteAddr().String())
	go func() {
		sum := 0
		for c := range requestTodo {
			sum += c
			if sum == 0 && closedRead {
				conn.(*net.TCPConn).CloseWrite()
				log.Printf("%s:close write", conn.RemoteAddr().String())
			}
		}
	}()
	for input.Scan() {
		content := input.Text()
		requestTodo <- 1
		go func() {
			_, err := fmt.Fprintln(conn, content+"!")
			if err != nil {
				log.Println(err)
			}
			time.Sleep(5 * time.Second)
			_, err = fmt.Fprintln(conn, content+"!!")
			if err != nil {
				log.Println(err)
			}
			time.Sleep(5 * time.Second)
			_, err = fmt.Fprintln(conn, content+"!!!")
			if err != nil {
				log.Println(err)
			}
			requestTodo <- -1
		}()
	}
	log.Printf("%s:close read", conn.RemoteAddr().String())
	conn.(*net.TCPConn).CloseRead()
	closedRead = true
	requestTodo <- 0
}

func main() {
	listenner, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err.Error())
	}
	for {
		conn, err := listenner.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}
		go handleConn(conn)
	}
}
