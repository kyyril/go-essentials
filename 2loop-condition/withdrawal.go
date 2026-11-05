package main

func CanWithdrawal(balance, amount, dailyLimit int) bool {
	if amount <= 0 {
		return false
	}
	if amount > balance {
		return false
	}
	if amount > dailyLimit {
		return false
	}
	return true
}
