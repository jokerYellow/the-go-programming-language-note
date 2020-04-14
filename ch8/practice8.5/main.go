package main

import (
	"image"
	"image/color"
	"math/cmplx"
	"sync"
	"time"
)

type point struct {
	width, py int
}

func main() {

}

func slow() {
	count := 10
	for i := 0; i < count; i++ {
		time.Sleep(time.Millisecond * 100)
	}
}

func fast() {
	poolCount := 3
	jobs := make(chan int, poolCount)
	count := 10
	wg := sync.WaitGroup{}
	for i := 0; i < poolCount; i++ {
		go func() {
			for range jobs {
				time.Sleep(time.Millisecond * 100)
				wg.Done()
			}
		}()
	}
	for i := 0; i < count; i++ {
		jobs <- i
		wg.Add(1)
	}
	wg.Wait()
}

func makePngFaster(width, height, poolsCount int) {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
	)
	wg := &sync.WaitGroup{}
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	jobs := make(chan point, poolsCount*2)
	for i := 0; i < poolsCount; i++ {
		go func() {
			for p := range jobs {
				for px := 0; px < p.width; px++ {
					y := float64(p.py)/float64(height)*(ymax-ymin) + ymin
					x := float64(px)/float64(width)*(xmax-xmin) + xmin
					z := complex(x, y)
					img.Set(px, p.py, mandelbrot(z))
				}
				wg.Done()
			}
		}()
	}
	for py := 0; py < height; py++ {
		wg.Add(1)
		jobs <- point{width, py}
	}
	wg.Wait()
}

func makePng(width, height int) {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
	)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/float64(height)*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/float64(width)*(xmax-xmin) + xmin
			z := complex(x, y)
			img.Set(px, py, mandelbrot(z))
		}
	}
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15
	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{255 - contrast*n}
		}
	}
	return color.Black
}
