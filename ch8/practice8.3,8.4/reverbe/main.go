package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func handleConn(conn net.Conn) {
	input := bufio.NewScanner(conn)
	wg := &sync.WaitGroup{}
	log.Printf("%s:connected", conn.RemoteAddr().String())
	for input.Scan() {
		wg.Add(1)
		go func(content string, delay time.Duration) {
			flag := ""
			for i := 0; i < 3; i++ {
				flag += "!"
				_, err := fmt.Fprintln(conn, content+flag)
				if err != nil {
					log.Println(err)
				}
				time.Sleep(delay)
			}
			wg.Done()
		}(input.Text(), 1*time.Second)
	}
	conn.(*net.TCPConn).CloseRead()
	log.Printf("%s:read closed", conn.RemoteAddr().String())
	wg.Wait()
	conn.(*net.TCPConn).CloseWrite()
	log.Printf("%s:write closed", conn.RemoteAddr().String())
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
