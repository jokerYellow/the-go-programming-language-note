package main

import (
	"fmt"
	"time"
)

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func fib(x int) int {
	if x <= 1 {
		return x
	}
	return fib(x-2) + fib(x-1)
}

func main() {
	go spinner(1 * time.Second)
	n := 44
	x := fib(n)
	fmt.Printf("\rFib(%d) is %d\n", n, x)
}
