package main

import "testing"

func TestCanWithdrawal(t *testing.T) {
	tests := []struct {
		name       string
		balance    int
		amount     int
		dailyLimit int
		want       bool
	}{
		{
			name:       "withdrawal with balance and limit",
			balance:    700,
			amount:     200,
			dailyLimit: 1000,
			want:       true,
		},
		{
			name:       "insufficient balance",
			balance:    300,
			amount:     500,
			dailyLimit: 1000,
			want:       false,
		},
		{
			name:       "exceeds daily limit",
			balance:    700,
			amount:     600,
			dailyLimit: 500,
			want:       false,
		},
		{
			name:       "negative withdrawal",
			balance:    700,
			amount:     -100,
			dailyLimit: 1000,
			want:       false,
		},
		{
			name:       "zero withdrawal",
			balance:    700,
			amount:     0,
			dailyLimit: 1000,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CanWithdrawal(
				tt.balance,
				tt.amount,
				tt.dailyLimit,
			)

			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
