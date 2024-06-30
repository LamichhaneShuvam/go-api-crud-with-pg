package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDatabase() {
	var err error

	DB, err = sql.Open("postgres", os.Getenv("DB_CONNECTION_URL"))

	if err != nil {
		log.Fatal("Error while connecting with the database", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Error while pinging the database", err)
	}

	log.Println("Connected to database")
}

func CloseDatabaseConnection() {
	if err := DB.Close(); err != nil {
		log.Fatal("Error closing connection with the database", err)
	}

	log.Print("Closed connection with the database")
}


func RunMigration () {
	err := createTodoTable()
	if err != nil {
		log.Fatal("Error while migration!", err)
	}
	log.Println("Created todo table successfully!")
}

func createTodoTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS todo (
			id SERIAL PRIMARY KEY,
			title VARCHAR(100) NOT NULL,
			completed BOOLEAN DEFAULT false,
			created_at timestamp DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp
		)
	`
	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
