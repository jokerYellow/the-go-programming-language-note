package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	tcpConn := conn.(*net.TCPConn)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, tcpConn)
		log.Println("close read")
		tcpConn.CloseRead()
		done <- struct{}{}
	}()
	_, err = io.Copy(tcpConn, os.Stdin)
	if err != nil {
		log.Println("ioerror:" + err.Error())
	}
	err = tcpConn.CloseWrite()
	log.Println("close write")
	if err != nil {
		log.Fatal(err)
	}
	<-done
	log.Println("program over")
}
