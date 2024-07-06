package models

import (
	"database/sql"
	"log"
	"time"
)

type User struct {
	ID          int       `json:"id"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	LastLoginAt time.Time `json:"last_login_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) Create(user *User) error {
	query := `
		INSERT INTO public.user (email, password)
		VALUES ($1, $2)
		RETURNING id, created_at
	`
	err := r.DB.QueryRow(query, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		log.Println("Error while creating user", err)
		return err
	}

	log.Print("Inserted user successfully!")
	return nil
}

func (r *UserRepository) GetById(id int) (*User, error) {
	query := `
		SELECT 
			id, email, password, created_at, 
			COALESCE(last_login_at, $2), COALESCE(updated_at, $2) 
			FROM public.user 
		WHERE id = $1
	`
	user := &User{}

	defaultTimeValue := time.Time{}
	err := r.DB.QueryRow(query, id, defaultTimeValue).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.LastLoginAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
	query := `
		SELECT 
			id, email, password, COALESCE(last_login_at, $2), created_at, COALESCE(updated_at, $2)
			FROM public.user
		WHERE email = $1
	`

	user := &User{}

	defaultDateValue := time.Time{}

	err := r.DB.QueryRow(query, email, defaultDateValue).Scan(&user.ID, &user.Email, &user.Password, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) DeleteById(id int) (*User, error) {
	query := `
		DELETE FROM user 
		WHERE id = $1
		RETURNING id, email, COALESCE(last_login_at, $2), created_at, COALESCE(updated_at, $2)
	`

	user := &User{}

	defaultDateValue := time.Time{}
	err := r.DB.QueryRow(query, id, defaultDateValue).Scan(&user.ID, &user.Email, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) UpdateUserPassword(id int, password string) error {
	query := `
		UPDATE public.user 
		SET password = $2, updated_at = current_timestamp
		WHERE id = $1
	`

	_, err := r.DB.Exec(query, id, password)

	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UpdateLastLoginAt(id int) error {
	query := `
		UPDATE public.user 
		SET last_login_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.DB.Exec(query, id)

	if err != nil {
		return err
	}
	return nil
}
