package models
import "time"

type User struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`// "-" agar tidak muncul saat di-convert ke JSON
	Role string `json:"role"`
}

type Task struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	UserId int `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}