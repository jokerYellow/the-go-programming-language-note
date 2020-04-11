package main

import (
	"fmt"
)

func main() {
	done := make(chan struct{})
	go func() {
		fmt.Println("done")
		done <- struct{}{}
		for i := 0; ; i++ {
			fmt.Printf("done%d\n", i)
		}
	}()
	<-done
}
