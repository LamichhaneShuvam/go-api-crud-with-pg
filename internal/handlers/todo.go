package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/lamichhaneshuvam/todo-pg/api"
	"github.com/lamichhaneshuvam/todo-pg/internal/db"
	"github.com/lamichhaneshuvam/todo-pg/internal/models"
)

func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	todoRepository := models.TodoRepository{DB: db.DB}

	if err := todoRepository.Create(&todo); err != nil {
		log.Println(err)
		api.InternalErrorHandler(w)
		return
	}

	api.CreateResponseHandler(w, todo, "Created todo successfully!")
	return
}

func GetTodoByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	todoRepository := models.TodoRepository{DB: db.DB}

	todo, err := todoRepository.GetById(id)

	if err != nil {
		if err == sql.ErrNoRows {
			api.NotFoundErrorHandler(w, errors.New("Todo not found!"))
			return
		}
		log.Println(err)
		api.InternalErrorHandler(w)
		return
	}

	api.OkResponseHandler(w, todo, "Fetched todo successfully!")
	return
}

func DeleteTodoByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		api.RequestErrorHandler(w, err)
	}

	todoRepository := models.TodoRepository{DB: db.DB}

	todo, err := todoRepository.DeleteById(id)

	if err != nil {
		if err == sql.ErrNoRows {
			api.NotFoundErrorHandler(w, errors.New("Todo not found!"))
			return
		}
		log.Println(err)
		api.InternalErrorHandler(w)
		return
	}

	api.OkResponseHandler(w, todo, "Deleted todo succssfully!")
	return
}

func UpdateTodoByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		api.RequestErrorHandler(w, err)
	}

	var todo models.Todo

	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	todoRepository := models.TodoRepository{DB: db.DB}

	//* find first and replace the null values
	oldDocs, err := todoRepository.GetById(id)

	if err != nil {
		if err == sql.ErrNoRows {
			api.NotFoundErrorHandler(w, errors.New("Todo not found!"))
			return
		}
		log.Println(err)
		api.InternalErrorHandler(w)
		return
	}

	if todo.Title == "" || todo.Title == " " {
		todo.Title = oldDocs.Title
	}

	//* Update the entry
	err = todoRepository.UpdateById(id, &todo)
	todo.UpdatedAt = time.Now()
	todo.ID = id

	if err != nil {
		log.Println(err)
		api.InternalErrorHandler(w)
		return
	}

	api.OkResponseHandler(w, todo, "Updated todo successfully!")
	return
}
