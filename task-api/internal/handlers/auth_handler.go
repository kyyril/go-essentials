package handlers
import (
	"task-api/internal/auth"
	"task-api/internal/database"
	"task-api/internal/models"
	"net/http"
	"encoding/json"
	"database/sql"
)

  
func Register(w http.ResponseWriter, r *http.Request){
	var user models.User
	//hash pw sebelum di simpan
	hashedPassword, _ := auth.HashPassword(user.Password)
	_, err := database.DB.Exec("INSERT INTO users (username, password, role) VALUES (?,?,?)",
								user.Username, hashedPassword, user.Role)
	if err != nil {
		http.Error(w, "username sudah ada", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User berhasil didaftarkan"})
}

func Login(w http.ResponseWriter, r *http.Request){
	var input models.User
	json(NewEncoder(r.body).Decode(input))
	var dbUser models.User

	err := database.Db.QueryRow("SELECT password, role FROM users WHERE username = ?", input.Username)
	if err == sql.ErrNoRows || !auth.CheckPassword(input.Password, dbUser.Password){
		http.Error(w, "invalid username or password", http.statusUnauthorized)
		return
	}
	token, _ := auth.GenerateJwt(input.Username, dbUser.Role)
	json.NewEncoder(w).Encode(map[string]string {"token": token})
}