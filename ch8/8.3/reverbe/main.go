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
	closedWrite := make(chan struct{})
	requestDone := make(chan struct{})
	log.Printf("%s:connected", conn.RemoteAddr().String())
	go func() {
		sum := 0
		for c := range requestTodo {
			sum += c
			if sum == 0 && (c == 0 || c == -1) {
				go func() {
					requestDone <- struct{}{}
				}()
			}
		}
	}()
	go func() {
		<-closedWrite
		<-requestDone
		conn.(*net.TCPConn).CloseWrite()
		log.Println("close write")
	}()
	for input.Scan() {
		log.Println("new scan")
		content := input.Text()
		requestTodo <- 1
		log.Println("new scan:", content)
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
	log.Println("goto close read")
	closedWrite <- struct{}{}
	log.Println("close read")
	requestTodo <- 0
	conn.(*net.TCPConn).CloseRead()
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
