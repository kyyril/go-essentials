package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite" // Import driver dengan blank identifier
)

var DB *sql.DB

func InitDB(){
	var err error
	DB, err := sql.Open("sqlite", "./data/task.db")
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	userTable := `
	CREATE TABLE IF NOT EXIST users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE,
			password TEXT,
			role TEXT
			)
	`

	taskTable := `
	CREATE TABLE IF NOT EXIST task (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			description TEXT,
			user_id INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP),
			FOREINGN KEY(user_id) REFERENCE users(id
			)`
	_, err = DB.Exec(userTable)
	if err != nil {
		log.Fatal("Gagal membuat tabel users:", err)
	}

	_, err = DB.Exec(taskTable)
	if err != nil {
		log.Fatal("Gagal membuat tabel tasks:", err)
	}

	log.Println("Database & Tabel berhasil disiapkan.")
}

