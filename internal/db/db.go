package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "host=localhost port=5432 user=postgres password=pstgrs dbname=marketwatch sslmode=disable"
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("DB unreachable:", err)
	}

	log.Println("Connected to Postgres")
}
