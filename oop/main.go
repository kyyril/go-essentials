package main

import "fmt"

// 1. ENKAPSULASI (Encapsulation)
// Kita pakai huruf besar (Wallet) agar bisa diakses luar,
// tapi isi saldonya (isi) huruf kecil agar rahasia/private.
type Wallet struct {
	Username string
	balance int //private
}

// Method untuk menambah balance dengan validasi (Enkapsulasi)
func (d *Wallet) topUp (amount int) bool{
	if amount > 0 {
		d.balance += amount
		fmt.Printf("Topup %s success for amount %d", d.Username, amount)
		return true
	}
	fmt.Printf("need amount")
	return false
}


// 2. POLIMORFISME (Interface)
// Kita bikin "Kontrak" bernama Payment. 
// Apa pun yang punya fungsi Pay(), dia dianggap Payment.
type Payment interface {
	Pay(nominal int) bool
}
// Strategi Bayar 1: Pakai Saldo Wallet
func (d *Wallet) Pay(nominal int)bool {
	if d.balance >= nominal{
		d.balance -= nominal
		fmt.Printf("%s pay %d. remaining balance: %d", d.Username, nominal, d.balance)
		return true
	}
	fmt.Println("Insufficient balance")
	return false
}

// Strategi Bayar 2: Pakai Kartu
type CreditCard struct {
	cardNumber string
}
func (c CreditCard) Pay(nominal int) bool {
	fmt.Printf("success pay %d with card: %s",nominal, c.cardNumber)
	return true
}

// 3. PEWARISAN/KOMPOSISI (Inherithence/Composition)
// PremiumMember "mewarisi" sifat Wallet lewat Embedding
type PremiumMember struct {
	Wallet 
	level string
}


func main (){
//create basic object wallet
budiWallet := Wallet{Username:"budi"}
budiWallet.topUp(1000)

//create with credit card
cardWahyu := CreditCard{cardNumber:"S8E-E21"}

//create object PremiumMember (Composition)
ucokMember := PremiumMember{
	Wallet: Wallet{Username:"ucok", balance:5000},
	level : "VIP",
}


//polimorphism: this function can pay with multi payment instrument
processTransaction := func(inst Payment, bil int){
	success := inst.Pay(bil)
	if success {
		fmt.Println("Payment successfull. Thanks!")
	}
	fmt.Println("-------------------------------")
}
//process transaction
processTransaction(&budiWallet, 500) //with wallet
processTransaction(cardWahyu, 2000) //card
processTransaction(&ucokMember, 1000) //member
}

// Encapsulation: Variabel balance tidak bisa diutak-atik langsung dari luar (misal wallet.balance = 0 itu dilarang jika beda package). Harus lewat fungsi TopUp.

// Inheritance/Composition: PremiumMember otomatis punya nama Pemilik dan bisa Bayar karena kita "menempelkan" Wallet di dalamnya. Tidak perlu menulis ulang kodenya.

// Polymorphism: Fungsi processTransaction tidak peduli kita bayar pakai kartu atau saldo, selama objek tersebut punya fungsi Pay(), program tetap jalan.