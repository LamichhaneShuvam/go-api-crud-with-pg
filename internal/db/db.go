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

func RunMigration() {
	var err error
	err = createTodoTable()
	if err != nil {
		log.Fatal("Error while running todo table migration!", err)
	}
	log.Println("Created todo table successfully!")

	err = createUserTable()
	if err != nil {
		log.Fatal("Error while running user table migration", err)
	}
	log.Println("Created user table successfully!")

	log.Println("Migration ran successfully!")
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

func createUserTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS public.user (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			last_login_at timestamp,
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
