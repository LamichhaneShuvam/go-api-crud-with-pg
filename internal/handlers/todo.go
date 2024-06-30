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
	"github.com/lamichhaneshuvam/todo-pg/internal/db"
	"github.com/lamichhaneshuvam/todo-pg/internal/models"
	"github.com/lamichhaneshuvam/todo-pg/internal/utils"
)

func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		utils.RequestErrorHandler(w, err)
		return
	}

	todoRepository := models.TodoRepository{DB: db.DB}

	if err := todoRepository.Create(&todo); err != nil {
		log.Println(err)
		utils.InternalErrorHandler(w)
		return
	}

	utils.CreateResponseHandler(w, todo, "Created todo successfully!")
	return
}

func GetTodoByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		utils.RequestErrorHandler(w, err)
		return
	}

	todoRepository := models.TodoRepository{DB: db.DB}

	todo, err := todoRepository.GetById(id)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.NotFoundErrorHandler(w, errors.New("Todo not found!"))
			return
		}
		log.Println(err)
		utils.InternalErrorHandler(w)
		return
	}

	utils.OkResponseHandler(w, todo, "Fetched todo successfully!")
	return
}

func DeleteTodoByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		utils.RequestErrorHandler(w, err)
	}

	todoRepository := models.TodoRepository{DB: db.DB}

	todo, err := todoRepository.DeleteById(id)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.NotFoundErrorHandler(w, errors.New("Todo not found!"))
			return
		}
		log.Println(err)
		utils.InternalErrorHandler(w)
		return
	}

	utils.OkResponseHandler(w, todo, "Deleted todo succssfully!")
	return
}

func UpdateTodoByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		utils.RequestErrorHandler(w, err)
	}

	var todo models.Todo

	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		utils.RequestErrorHandler(w, err)
		return
	}

	todoRepository := models.TodoRepository{DB: db.DB}

	//* find first and replace the null values
	oldDocs, err := todoRepository.GetById(id)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.NotFoundErrorHandler(w, errors.New("Todo not found!"))
			return
		}
		log.Println(err)
		utils.InternalErrorHandler(w)
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
		utils.InternalErrorHandler(w)
		return
	}

	utils.OkResponseHandler(w, todo, "Updated todo successfully!")
	return
}
