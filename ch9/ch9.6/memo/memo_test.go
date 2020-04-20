package memo

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

func httpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func Test(t *testing.T) {
	m := New(httpGetBody)
	incomes := []string{
		"https://www.baidu.com",
		"https://icon.wuruihong.com/",
		"https://www.zhihu.com/",
		"https://icon.wuruihong.com/",
		"https://www.zhihu.com/",
		"https://www.zhihu.com/",
		"https://www.zhihu.com/",
		"https://www.zhihu.com/",
		"https://www.zhihu.com/",
	}
	wg := &sync.WaitGroup{}
	s := time.Now()
	for _, url := range incomes {
		wg.Add(1)
		go func(url string) {
			start := time.Now()
			o, err := m.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			duration := time.Since(start)
			length := 0
			if o != nil {
				length = len(o.([]byte))
			}
			fmt.Printf("%s, %s, %d bytes\n", url, duration, length)
			wg.Done()
		}(url)
	}
	wg.Wait()
	fmt.Printf("sum duration: %s", time.Since(s))
}
