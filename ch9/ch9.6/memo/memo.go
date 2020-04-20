package memo

import "sync"

type entry struct {
	result
	ready chan struct{}
}

type Memo struct {
	f     Func
	cache map[string]*entry
	mu    *sync.Mutex
}

type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

func New(f Func) *Memo {
	return &Memo{f, make(map[string]*entry), new(sync.Mutex)}
}

func (memo *Memo) Get(key string) (interface{}, error) {
	memo.mu.Lock()
	res, ok := memo.cache[key]
	if !ok {
		res = new(entry)
		res.ready = make(chan struct{})
		memo.cache[key] = res
		memo.mu.Unlock()
		go memo.get(res, key)
	} else {
		memo.mu.Unlock()
	}
	<-res.ready
	return res.value, res.err
}

func (memo *Memo) get(e *entry, key string) {
	e.value, e.err = memo.f(key)
	close(e.ready)
}
