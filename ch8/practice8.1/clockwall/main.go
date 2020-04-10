package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type timeItem struct {
	c          chan string
	name, time string
}

const moveUp = "\x1b[1A"
const clearEntireLine = "\x1b[2k"

func handleConn(c net.Conn, ch chan string) {
	defer c.Close()
	for {
		b := make([]byte, 100)
		n, err := c.Read(b)
		if err != nil {
			log.Println(err)
			return
		}
		ch <- string(b[:n-1])
		time.Sleep(1 * time.Second)
	}
}

func (t *timeItem) listen() {
	for n := range t.c {
		t.time = n
	}
}

func main() {
	items := []*timeItem{}
	for _, r := range os.Args[1:] {
		params := strings.Split(r, "=")
		tz := params[0]
		address := params[1]
		conn, err := net.Dial("tcp", address)
		if err != nil {
			log.Println(err)
			continue
		}
		c := make(chan string)
		go handleConn(conn, c)
		item := timeItem{c, tz, ""}
		go item.listen()
		items = append(items, &item)
	}
	for {
		for _, i := range items {
			fmt.Printf("%5s:\t%s\n", i.name, i.time)
		}
		time.Sleep(1 * time.Second)
		for i := 0; i < len(items); i++ {
			fmt.Print(moveUp)
			fmt.Print(clearEntireLine)
		}
	}
}
