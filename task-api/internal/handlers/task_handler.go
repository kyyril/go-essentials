package handlers

import (
	"encoding/json"
	"net/http"
	"task-api/internal/database"
	"task-api/internal/models"
)

// CreateTask - Menambahkan tugas baru
func CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Simpan ke database
	_, err := database.DB.Exec("INSERT INTO tasks (title, description, user_id) VALUES (?, ?, ?)",
		task.Title, task.Description, task.UserID)
	
	if err != nil {
		http.Error(w, "Gagal menyimpan task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Task created successfully"})
}

// GetTasks - Mengambil semua tugas
func GetTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, title, description, user_id, created_at FROM tasks")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		rows.Scan(&t.ID, &t.Title, &t.Description, &t.UserID, &t.CreatedAt)
		tasks = append(tasks, t)
	}

	json.NewEncoder(w).Encode(tasks)
}

// DeleteTask - Menghapus tugas (Hanya Admin atau Pemilik)
func DeleteTask(w http.ResponseWriter, r *http.Request) {
    // Di sini logika role check akan bekerja via middleware
    id := r.URL.Query().Get("id")
    if id == "" {
        http.Error(w, "ID required", http.StatusBadRequest)
        return
    }

    _, err := database.DB.Exec("DELETE FROM tasks WHERE id = ?", id)
    if err != nil {
        http.Error(w, "Gagal menghapus task", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted"})
}