package bank

import "sync"

var balance int

var mu = sync.RWMutex{}

func Deposit1(amount int) {
	mu.Lock()
	defer func() {
		mu.Unlock()
	}()
	balance += amount
}

func Balance1() int {
	mu.RLock()
	defer func() {
		mu.RUnlock()
	}()
	return balance
}
