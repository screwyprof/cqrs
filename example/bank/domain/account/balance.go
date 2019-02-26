package account

type Balance struct {
	amount uint64
}

func NewBalance() Balance {
	return Balance{}
}

func (b Balance) Deposit(amount uint64) Balance {
	return Balance{
		amount: b.amount + amount,
	}
}

func (b Balance) Withdraw(amount uint64) Balance {
	return Balance{
		amount: b.amount - amount,
	}
}

func (b Balance) WithdrawlWillResultInNegativeBalance(amount uint64) bool {
	return (b.amount - amount) < 0
}
