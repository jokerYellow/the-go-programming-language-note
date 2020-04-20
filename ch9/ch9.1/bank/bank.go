package bank

import (
	"fmt"
)

type Result struct {
	rt      bool
	balance int
}

type withDrawItem struct {
	amount int
	rt     chan Result
}

var balances = make(chan int)
var withDraw = make(chan withDrawItem)

func Deposit(amount int) Result {
	return WithDraw(-amount)
}

func Balance() int {
	return <-balances
}

func WithDraw(amount int) Result {
	rt := make(chan Result)
	withDraw <- withDrawItem{amount, rt}
	return <-rt
}

func teller() int {
	var balance int
	for {
		select {
		case item := <-withDraw:
			if balance-item.amount < 0 {
				item.rt <- Result{false, balance}
				fmt.Printf("bank withdraw failed:%d balance:%d\n", item.amount, balance)
				continue
			}
			before := balance
			balance -= item.amount
			item.rt <- Result{true, balance}
			fmt.Printf("bank change before:%d success:%d balance:%d\n", before, -item.amount, balance)
		case balances <- balance:
			fmt.Printf("bank balance:%d\n", balance)
		}
	}
}

func init() {
	go teller()
}
