package main

import (
	"fmt"
	"testing"
)

const width = 1024
const height = 1024

func init() {
	fmt.Println("2.5 GHz 四核Intel Core i7")
}

func Benchmark(b *testing.B) {
	b.Run("slow", func(b *testing.B) {
		makePng(width, height)
	})
	for i := 1; i < 40; i++ {
		b.Run(fmt.Sprint(i), func(b *testing.B) {
			makePngFaster(width, height, i)
		})
	}
}

func TestFast(t *testing.T) {
	fast()
}

func TestSlow(t *testing.T) {
	slow()
}
