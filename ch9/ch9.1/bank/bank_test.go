package bank

import (
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T) {
	for i := 0; i < 100; i++ {
		go Deposit(100)
		go Deposit(200)
	}
	time.Sleep(1 * time.Second)
	fmt.Printf("balance:%d\n", Balance())
	if Balance() != 30000 {
		t.Fail()
	}
}

func Test1(t *testing.T) {
	deposit := func(count int, index int, step int) {
		fmt.Printf("before balance:%4d amount:%4d index:%d %d\n", Balance1(), count, index, step)
		Deposit1(count)
		fmt.Printf("after  balance:%4d amount:%4d index:%d %d\n", Balance1(), count, index, step)
	}
	for i := 0; i < 100; i++ {
		go deposit(100, i, 0)
		go deposit(200, i, 1)
	}
	time.Sleep(3 * time.Second)
	if Balance1() != 30000 {
		t.Fail()
	}
}

func TestNotSafe(t *testing.T) {
	deposit := func(count int) {
		rt := DepositNotSafe(count)
		fmt.Printf("before:%d balance:%d balance:%t amount:%d \n", rt.before, rt.balance, rt.rt, count)
	}
	for i := 0; i < 100; i++ {
		go deposit(100)
		go deposit(200)
	}
	time.Sleep(1 * time.Second)
	if Balance() != 30000 {
		t.Fail()
	}
	fmt.Printf("balance:%d\n", BalanceNotSafe())
}
