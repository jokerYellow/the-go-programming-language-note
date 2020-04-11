package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
)

var cmdsMap = map[string]func(string, net.Conn){}

func init() {
	cmdsMap["cd"] = cdCmd
	cmdsMap["ls"] = lsCmd
	cmdsMap["ll"] = lsCmd
	cmdsMap["pwd"] = pwdCmd
	cmdsMap["cat"] = catCmd
	cmdsMap["get"] = getCmd
}

func command(cmd string, conn net.Conn) {
	if len(cmd) == 0 {
		return
	}
	cmds := strings.Split(cmd, " ")
	if c, ok := cmdsMap[cmds[0]]; ok {
		c(cmd, conn)
		return
	}
	defaultCmd(cmd, conn)
}

func handleConn(conn net.Conn) {
	log.Printf("%s: connected", conn.RemoteAddr().String())
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		cmd := scanner.Text()
		log.Printf("%s: %s\n", conn.RemoteAddr().String(), cmd)
		command(cmd, conn)
	}
	defer func() {
		conn.Close()
		log.Printf("%s: closed", conn.RemoteAddr().String())
	}()
}

func main() {
	port := os.Args[1]
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ftp server start at localhost:%s\n", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}
