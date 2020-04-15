package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	flag.Parse()
	go func() {
		os.Stdin.Read(make([]byte, 1))
		close(done)
	}()
	du(flag.Args())
}

var done = make(chan struct{})

func cancel() bool {
	select {
	case <-done: //当done关闭后，会一直触发该方法，其实就是_,ok:= <- done，ok=false
		return true
	default:
		return false
	}
}

func du(roots []string) {
	fileSizes := make(chan int64)
	tick := time.NewTicker(500 * time.Millisecond)
	wg := &sync.WaitGroup{}
	for _, dir := range roots {
		wg.Add(1)
		go workDir(dir, fileSizes, wg)
	}
	go func() {
		wg.Wait()
		close(fileSizes)
	}()
	var bytes int64 = 0
	fileNums := 0

loop:
	for {
		select {
		case <-done:
			for range fileSizes {
			}
			return
		case fileSize, ok := <-fileSizes:
			if !ok {
				break loop
			}
			bytes += fileSize
			fileNums++
		case <-tick.C:
			printDiskUsage(fileNums, bytes)
		}
	}
	tick.Stop()
	printDiskUsage(fileNums, bytes)
	if cancel() {
		panic("main should not end normally")
	}
}

func printDiskUsage(nfiles int, bytes int64) {
	fmt.Printf("%d files  %.1fMBytes\n", nfiles, float64(bytes)/1e6)
}

var sema = make(chan struct{}, 30)

func workDir(dir string, fileSizes chan<- int64, wg *sync.WaitGroup) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			if wg != nil {
				wg.Add(1)
			}
			go workDir(subdir, fileSizes, wg)
		} else {
			fileSizes <- entry.Size()
		}
	}
	wg.Done()
}

func dirents(dir string) []os.FileInfo {
	if cancel() {
		return nil
	}
	sema <- struct{}{}
	defer func() { <-sema }()
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return nil
	}
	return entries
}
