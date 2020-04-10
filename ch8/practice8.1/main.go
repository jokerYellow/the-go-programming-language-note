package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func handleConn(c net.Conn, delay time.Duration) {
	for {
		_, err := io.WriteString(c, time.Now().Format("2006 Mon Jan 2 15:04:05 -0700 MST\n"))
		if err != nil {
			log.Println(err)
			return
		}
		time.Sleep(delay)
	}
}

func main() {
	port := os.Args[1]
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn, 1*time.Second)
	}
}
