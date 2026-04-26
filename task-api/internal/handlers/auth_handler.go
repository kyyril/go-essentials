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
	json.NewDecoder(r.Body).Decode(&user)
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
	json.NewDecoder(r.Body).Decode(&input)
	var dbUser models.User

	row := database.DB.QueryRow("SELECT password, role FROM users WHERE username = ?", input.Username)
	err := row.Scan(&dbUser.Password, &dbUser.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "invalid username or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if !auth.CheckPassword(input.Password, dbUser.Password) {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}
	token, _ := auth.GenerateJwt(input.Username, dbUser.Role)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}