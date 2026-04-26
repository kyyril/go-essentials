package auth

import "testing"

func TestHashPassword(t *testing.T) {
	password := "rahasia123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("Gagal hashing: %v", err)
	}

	if !CheckPassword(password, hash) {
		t.Errorf("Password tidak cocok dengan hash")
	}
}