package memo

import (
	"fmt"
	"time"
)

type entry struct {
	result
	ready chan struct{}
}

type request struct {
	key      string
	response chan result
}

type Memo struct {
	f       Func
	cache   map[string]*entry
	request chan request
}

type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

func New(f Func) *Memo {
	memo := Memo{f, make(map[string]*entry), make(chan request)}
	go memo.server()
	return &memo
}

func (memo *Memo) Get(key string) (interface{}, error) {
	response := make(chan result)
	memo.request <- request{key, response}
	r := <-response
	return r.value, r.err
}

func (memo *Memo) server() {
	for r := range memo.request {
		e := memo.cache[r.key]
		if e == nil {
			e = &entry{ready: make(chan struct{})}
			memo.cache[r.key] = e
			go func(item request, e *entry) {
				e.value, e.err = memo.f(item.key)
				fmt.Printf("%s : %s\n", item.key, time.Now())
				close(e.ready)
			}(r, e)
		}
		go e.deliver(r.response)
	}
}

func (memo *Memo) Close() {
	close(memo.request)
}

func (e *entry) deliver(response chan<- result) {
	<-e.ready
	response <- e.result
}
