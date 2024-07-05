// cmd/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lamichhaneshuvam/todo-pg/internal/db"
	"github.com/lamichhaneshuvam/todo-pg/internal/handlers"
	"github.com/lamichhaneshuvam/todo-pg/internal/middlewares"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	//* Load the environment variables
	godotenv.Load()

	//* Database connection
	db.InitDatabase()

	//* Database migration
	db.RunMigration()

	//* Close connection with the db on defer
	defer db.CloseDatabaseConnection()

	router := mux.NewRouter()

	//* Routes
	// [TODO ROUTES]
	router.HandleFunc("/todos", handlers.CreateTodoHandler).Methods("POST")
	router.HandleFunc("/todos/{id:[0-9]+}", middlewares.UserAuth(handlers.GetTodoByIdHandler)).Methods("GET")
	router.HandleFunc("/todos/{id:[0-9]+}", handlers.DeleteTodoByIdHandler).Methods("DELETE")
	router.HandleFunc("/todos/{id:[0-9]+}", handlers.UpdateTodoByIdHandler).Methods("PUT")

	// [USER ROUTES]
	router.HandleFunc("/user", handlers.CreateUserHandler).Methods("POST")
	router.HandleFunc("/user/login", handlers.LoginUserHandler).Methods("POST")

	//* Server starts
	var PORT string = os.Getenv("APPLICATION_PORT")
	addr := fmt.Sprintf(":%s", PORT)

	log.Println("Go Server Started on", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
