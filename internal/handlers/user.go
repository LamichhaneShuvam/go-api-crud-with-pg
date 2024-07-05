package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"strings"

	"net/http"

	"github.com/lamichhaneshuvam/todo-pg/internal/db"
	"github.com/lamichhaneshuvam/todo-pg/internal/models"
	"github.com/lamichhaneshuvam/todo-pg/internal/utils"
)

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.RequestErrorHandler(w, err)
		return
	}
	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		log.Println(err)
		utils.InternalErrorHandler(w)
		return
	}

	formattedEmail := strings.ToLower(strings.Trim(user.Email, " "))
	user.Email = formattedEmail

	userRepository := models.UserRepository{DB: db.DB}

	//* Check if the user already exists
	userExists, err := userRepository.GetByEmail(user.Email)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
			utils.InternalErrorHandler(w)
			return
		}
	}
	if userExists.Email != "" {
		utils.ConflictErrorHandler(w, errors.New("User with the same email already exists!"))
		return
	}

	user.Password = hashedPassword

	if err := userRepository.Create(&user); err != nil {
		log.Println(err)
		utils.InternalErrorHandler(w)
		return
	}

	utils.CreateResponseHandler(w, user, "Created user successfully, please log in to continue")
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var response loginResponse

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.RequestErrorHandler(w, err)
		return
	}

	formattedEmail := strings.ToLower(strings.Trim(user.Email, " "))
	user.Email = formattedEmail

	userRepository := models.UserRepository{DB: db.DB}

	userDetails, err := userRepository.GetByEmail(user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.NotFoundErrorHandler(w, errors.New("email or password incorrect!"))
			return
		}
		log.Println(err)
		utils.InternalErrorHandler(w)
		return
	}

	isPasswordValid := utils.CheckPasswordHash(user.Password, userDetails.Password)

	if !isPasswordValid {
		utils.NotFoundErrorHandler(w, errors.New("email or password incorrect!"))
		return
	}

	jwtToken, err := utils.GenerateJwt(userDetails.ID)

	if err != nil {
		log.Println(err)
		utils.InternalErrorHandler(w)
		return
	}

	response.AccessToken = jwtToken
	response.RefreshToken = jwtToken

	utils.OkResponseHandler(w, response, "User logged in successfully")
	return
}
