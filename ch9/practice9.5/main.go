package main

import (
	"fmt"
	"time"
)

//写一个程序，两个 goroutine 通过两个无缓冲通道来相互转发消息。这个程序能维持美妙多少次通信？
/*
1 second count: 1960425
1 second count: 1916263
1 second count: 1797235
1 second count: 1825926
1 second count: 1833524
1 second count: 1835546
*/
func main() {
	ch1 := make(chan struct{})
	ch2 := make(chan struct{})
	go func() {
		for {
			c := <-ch1
			ch2 <- c
		}
	}()
	ch1 <- struct{}{}
	ticker := time.NewTicker(1 * time.Second)
	count := 0
	for {
		select {
		case c := <-ch2:
			count++
			ch1 <- c
		case <-ticker.C:
			fmt.Printf("1 second count: %d\n", count)
			count = 0
		}
	}
}
