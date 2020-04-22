package main

import (
	"fmt"
	"time"
)

//使用通道构造一个把任意多个goroutine串联在一起的流水线程序，在内存耗尽之前你能创建的最大流水线级数是多少？一个值穿过整个流水线需要多久？
/*
输出结果
count:1000000 duration:425.254693ms
count:1100000 duration:553.257326ms
count:1200000 duration:974.551749ms
count:1300000 duration:802.393849ms
count:1400000 duration:901.129533ms
count:1500000 duration:1.012876386s
count:1600000 duration:1.1125231s
...
count:4700000 duration:22.359718s
count:4800000 duration:24.323349042s
count:4900000 duration:19.476047948s
count:5000000 duration:25.533480372s
*/
func main() {
	count := 1000000
	step := 100000
	for {
		makePipline(count)
		count += step
	}
}

func makePipline(count int) {
	ches := make([]chan struct{}, count)
	for i := 0; i < count; i++ {
		ches[i] = make(chan struct{})
		if i > 0 {
			go pip(ches[i-1], ches[i])
		}
	}
	var start time.Time
	go func() {
		start = time.Now()
		close(ches[0])
	}()
	<-ches[count-1]
	fmt.Printf("count:%d duration:%s\n", count, time.Since(start))
}

func pip(in <-chan struct{}, out chan<- struct{}) {
	c := <-in
	out <- c
}
