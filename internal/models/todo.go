package models

import (
	"database/sql"
	"log"
	"time"
)

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TodoRepository struct {
	DB *sql.DB
}

func (r *TodoRepository) Create(todo *Todo) error {
	query := `
		INSERT INTO todo (title, completed)
		VALUES ($1, $2)
		RETURNING id, created_at
	`
	err := r.DB.QueryRow(query, todo.Title, todo.Completed).Scan(&todo.ID, &todo.CreatedAt)

	if err != nil {
		log.Println("Error while creating todo", err)
		return err
	}

	log.Print("Inserted todo successfully!")
	return nil
}

func (r *TodoRepository) GetById(id int) (*Todo, error) {
	query := `
		SELECT id, title, completed, created_at, COALESCE(updated_at, $2) FROM todo
		WHERE id = $1
	`
	todo := &Todo{}

	defaultValue := time.Time{}
	err := r.DB.QueryRow(query, id, defaultValue).Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (r *TodoRepository) DeleteById(id int) (*Todo, error) {
	query := `
		DELETE FROM todo
		WHERE id = $1
		RETURNING id, title, completed, created_at, COALESCE(updated_at, $2)
	`

	todo := &Todo{}

	defaultDateValue := time.Time{}
	err := r.DB.QueryRow(query, id, defaultDateValue).Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (r *TodoRepository) UpdateById(id int, todo *Todo) error {
	query := `
		UPDATE todo
		SET title = $2, completed = $3, updated_at = current_timestamp
		WHERE id = $1
	`

	_, err := r.DB.Exec(query, id, todo.Title, todo.Completed)

	if err != nil {
		return err
	}
	return nil
}
