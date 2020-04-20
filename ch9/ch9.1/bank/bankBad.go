package bank

var balanceNotSafe int

type ResultNotSafe struct {
	before, balance int
	rt              bool
}

func DepositNotSafe(amount int) ResultNotSafe {
	before := balanceNotSafe
	if balanceNotSafe+amount < 0 {
		return ResultNotSafe{before, balanceNotSafe, false}
	}
	balanceNotSafe += amount
	return ResultNotSafe{before, balanceNotSafe, true}
}

func BalanceNotSafe() int {
	return balanceNotSafe
}
