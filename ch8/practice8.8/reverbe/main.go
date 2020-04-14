package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const delay = 1 * time.Second
const timeout = 5

//使用select语句，给回声服务器加一个超时，断开10s内没有任何呼叫的客户端
func handleConn(conn net.Conn) {
	input := bufio.NewScanner(conn)
	ch := make(chan string)
	wg := &sync.WaitGroup{}
	log.Printf("%s: connected", conn.RemoteAddr().String())
	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		ticker.Stop()
		conn.Close()
		log.Printf("%s: closed", conn.RemoteAddr().String())
	}()
	done := make(chan struct{}, 1)
	go func() {
		countdown := timeout
		for {
			select {
			case content := <-ch:
				countdown = timeout
				wg.Add(1)
				go echo(content, conn, delay, wg)
			case <-ticker.C:
				countdown--
				log.Printf("%s: last %d seconds", conn.RemoteAddr().String(), countdown)
				fmt.Fprintf(conn, "timeout in %d seconds\n", countdown)
				if countdown == 0 {
					done <- struct{}{}
					return //超时之后，返回该函数，避免死循环
				}
			}
		}
	}()
	go func() {
		for input.Scan() {
			ch <- input.Text()
		}
		done <- struct{}{}
	}()
	<-done
	wg.Wait()
}

func echo(content string, conn net.Conn, delay time.Duration, wg *sync.WaitGroup) {
	flag := ""
	for i := 0; i < 3; i++ {
		flag += "!"
		_, err := fmt.Fprintln(conn, content+flag)
		log.Printf("%s:send %s", conn.RemoteAddr().String(), content+flag)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(delay)
	}
	wg.Done()
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
