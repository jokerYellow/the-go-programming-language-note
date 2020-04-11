package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err.Error())
	}
	go io.Copy(os.Stdout, conn)
	io.Copy(conn, os.Stdin)
	defer conn.Close()
}
