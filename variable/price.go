package main

const taxRate = 0.11

func CalculatePrice(price float64) float64 {
	taxPrice := price * taxRate
	return taxPrice
}
