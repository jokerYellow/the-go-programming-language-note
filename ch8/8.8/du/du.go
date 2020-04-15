package du

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func du(roots []string) {
	fileSizes := make(chan int64)
	go func() {
		for _, dir := range roots {
			workDir(dir, fileSizes)
		}
		close(fileSizes)
	}()
	var bytes int64 = 0
	fileNums := 0
	for fileSize := range fileSizes {
		bytes += fileSize
		fileNums++
	}
	printDiskUsage(fileNums, bytes)
}

func duShowProgress(roots []string) {
	fileSizes := make(chan int64)
	tick := time.NewTicker(500 * time.Millisecond)
	go func() {
		for _, dir := range roots {
			workDir(dir, fileSizes)
		}
		close(fileSizes)
	}()
	var bytes int64 = 0
	fileNums := 0

loop:
	for {
		select {
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
}

func DuFaster(roots []string) {
	fileSizes := make(chan int64)
	wg := &sync.WaitGroup{}
	for _, dir := range roots {
		wg.Add(1)
		go workDir2(dir, fileSizes, wg)
	}
	go func() {
		wg.Wait()
		close(fileSizes)
	}()
	var bytes int64 = 0
	fileNums := 0
	for fileSize := range fileSizes {
		bytes += fileSize
		fileNums++
	}
	printDiskUsage(fileNums, bytes)
}

func duShowProgressFaster(roots []string) {
	fileSizes := make(chan int64)
	tick := time.NewTicker(500 * time.Millisecond)
	wg := &sync.WaitGroup{}
	for _, dir := range roots {
		wg.Add(1)
		go workDir2(dir, fileSizes, wg)
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
}

func printDiskUsage(nfiles int, bytes int64) {
	fmt.Printf("%d files  %.1fMBytes\n", nfiles, float64(bytes)/1e6)
}

var sema = make(chan struct{}, 30)

func workDir(dir string, fileSizes chan<- int64) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			workDir(subdir, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

func workDir2(dir string, fileSizes chan<- int64, wg *sync.WaitGroup) {
	sema <- struct{}{}
	defer func() { <-sema }()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			if wg != nil {
				wg.Add(1)
			}
			go workDir2(subdir, fileSizes, wg)
		} else {
			fileSizes <- entry.Size()
		}
	}
	if wg != nil {
		wg.Done()
	}
}

func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return nil
	}
	return entries
}
