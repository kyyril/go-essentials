package main

const taxRate = 0.1 //10%

func CalculatePrice(price float64) float64 {
	return price + CalculateTax(price)
}

func CalculateTax(price float64) float64 {
	tax := price * taxRate
	return tax
}
