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
	for input.Scan() {
		content := input.Text()
		go func() {
			fmt.Fprintln(conn, content+"!")
			time.Sleep(1 * time.Second)
			fmt.Fprintln(conn, content+"!!")
			time.Sleep(1 * time.Second)
			fmt.Fprintln(conn, content+"!!!")
		}()
	}
	defer conn.Close()
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
		handleConn(conn)
	}
}
