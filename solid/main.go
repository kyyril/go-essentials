package main

import "fmt"

// 1. INTERFACE (Abstraksi - Prinsip DIP & ISP)
// Kita bikin interface kecil-kecil agar tidak "gemuk" (ISP).
type Notifier interface {
	Send(msg string) error
}

// 2. IMPLEMENTASI (Struct Konkret)
// EmailService fokus hanya pada Email (SRP)
type EmailService struct {
	EmailTujuan string
}

func (e EmailService) Send(msg string) error {
	fmt.Printf("Mengirim Email ke %s: %s\n", e.EmailTujuan, msg)
	return nil
}

// SMSService fokus hanya pada SMS (SRP)
type SMSService struct {
	NomorHP string
}

func (s SMSService) Send(msg string) error {
	fmt.Printf("Mengirim SMS ke %s: %s\n", s.NomorHP, msg)
	return nil
}

// 3. SERVICE UTAMA (High Level Module)
type ShippingService struct {
	// ShippingService tidak peduli mau kirim via Email atau SMS.
	// Dia hanya butuh sesuatu yang punya fungsi Send() (DIP).
	Pengirim Notifier
}

func (s ShippingService) UpdateStatus(orderID string, status string) {
	pesan := fmt.Sprintf("Order %s sekarang berstatus: %s", orderID, status)
	
	// Eksekusi pengiriman tanpa tahu ini Email atau SMS (Polimorfisme)
	err := s.Pengirim.Send(pesan)
	if err != nil {
		fmt.Println("Gagal mengirim notifikasi")
	}
}

func main() {
	// SKENARIO A: Kirim via Email
	notifEmail := EmailService{EmailTujuan: "user@example.com"}
	serviceA := ShippingService{Pengirim: notifEmail}
	serviceA.UpdateStatus("INV-001", "DIKIRIM")

	fmt.Println("-------------------------------")

	// SKENARIO B: Kirim via SMS (Tanpa mengubah kode ShippingService sama sekali - OCP)
	notifSMS := SMSService{NomorHP: "08123456789"}
	serviceB := ShippingService{Pengirim: notifSMS}
	serviceB.UpdateStatus("INV-002", "SAMPAI")
}

// S (Single Responsibility): EmailService cuma urus email, ShippingService cuma urus logika pengiriman.

// O (Open/Closed): Kalau besok bos minta tambah notifikasi WhatsApp, Anda cukup bikin struct WhatsAppService baru. Anda TIDAK PERLU menyentuh atau mengubah kode di dalam ShippingService.

// L (Liskov Substitution): EmailService dan SMSService bisa saling menggantikan posisi Notifier tanpa membuat program error.

// I (Interface Segregation): Interface Notifier sangat simpel, cuma ada satu fungsi Send. Tidak memaksa pengirim untuk punya fungsi yang tidak perlu (seperti Login atau Attachment).

// D (Dependency Inversion): ShippingService tidak "nempel" ke EmailService. Dia nempel ke interface Notifier. Ini rahasia supaya kode mudah di-unit test (Anda bisa bikin "Email Palsu/Mock" saat testing).