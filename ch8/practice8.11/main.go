package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

//当第一个请求返回的时候，取消其他的请求

var done = make(chan struct{})

func isCancel() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

func main() {
	fmt.Println(mirroredQuery())
	os.Stdin.Read(make([]byte, 1))
}

func mirroredQuery() string {
	responses := make(chan string, 3)
	go requestHost("1", responses)
	go requestHost("2", responses)
	go requestHost("3", responses)
	r, _ := <-responses
	fmt.Printf("%s: mirroredQuery\n", r)
	close(done)
	return r
}

func requestHost(hostname string, response chan<- string) {
	fmt.Printf("%s: begin requestHost\n", hostname)
	rt := make(chan string)
	go func() {
		r := request(hostname)
		fmt.Printf("%s: request done\n", hostname)
		rt <- r
	}()
loop:
	for {
		select {
		case <-done:
			fmt.Printf("%s: cancelled\n", hostname)
			break loop
		case h, _ := <-rt:
			if isCancel() {
				fmt.Printf("%s: append to chan but already cancelled\n", hostname)
				break loop
			}
			fmt.Printf("%s: insert to response\n", hostname)
			response <- h
			break loop
		}
	}
	fmt.Printf("%s: end requestHost\n", hostname)
}

func request(hostname string) string {
	a := time.Duration(rand.Intn(4)) //simulator request
	time.Sleep(a * time.Second)
	return hostname
}
