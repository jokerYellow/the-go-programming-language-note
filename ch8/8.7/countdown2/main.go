package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1))
		abort <- struct{}{}
	}()
	fmt.Println("Commencing countdown.Press return to abort.")
	ticker := time.NewTicker(1 * time.Second)
	for countdown := 10; countdown > 0; countdown-- {
		//select 后面跟的是chan的操作，如ch <-, <- ch。
		//每次select会选择非阻塞的ch来执行，如果有多个非阻塞的，则随机执行一个。
		//如果select后面为空 select {}，则会永久阻塞
		select {
		case <-ticker.C:
			fmt.Println(countdown)
		case <-abort:
			fmt.Println("abort launch.")
			return
		}
	}
	ticker.Stop()
	launch()
}

func launch() {
	fmt.Println("launch.")
}
