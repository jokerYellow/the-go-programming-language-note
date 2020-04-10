package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func command(cmd string) string {
	switch cmd {
	case "pwd":
		path, err := os.Getwd()
		if err != nil {
			return err.Error()
		}
		return path
	case "ls":
		path, err := os.Getwd()
		if err != nil {
			return err.Error()
		}
		files, _ := ioutil.ReadDir(path)
		b := bytes.Buffer{}
		for _, f := range files {
			b.WriteString(f.Name())
			b.WriteString("\n")
		}
		return b.String()
	}
}

func handleConn(conn net.Conn) {
	log.Printf("%s: connected", conn.RemoteAddr().String())
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		cmd := scanner.Text()
		log.Printf("%s: %s\n", conn.RemoteAddr().String(), cmd)
		conn.Write([]byte(command(cmd)))
	}
	defer func() {
		conn.Close()
		log.Printf("%s: closed", conn.RemoteAddr().String())
	}()
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}
