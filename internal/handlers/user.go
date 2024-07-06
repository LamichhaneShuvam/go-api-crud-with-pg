package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"

	"net/http"

	"github.com/lamichhaneshuvam/todo-pg/internal/db"
	"github.com/lamichhaneshuvam/todo-pg/internal/models"
	"github.com/lamichhaneshuvam/todo-pg/internal/utils"
)

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type PasswordChangeRequest struct {
	Password string `json:"password"`
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

	formattedEmail := strings.ToLower(strings.TrimSpace(user.Email))
	user.Email = formattedEmail

	user.Password = strings.TrimSpace(user.Password)

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

	if userExists != nil {
		utils.ConflictErrorHandler(w, errors.New("user with the same email already exists"))
		return
	}

	user.Password = hashedPassword

	if err := userRepository.Create(&user); err != nil {
		log.Println(err)
		utils.InternalErrorHandler(w)
		return
	}

	utils.CreateResponseHandler(w, user, "Created user successfully, please log in to continue")
	return
}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	var changePasswordPayload PasswordChangeRequest
	userId, err := strconv.Atoi(r.Header.Get("userId"))

	if err != nil {
		log.Println("Error while string to int conversion:", err)
		utils.UnauthorizedErrorHandler(w, errors.New("please login to continue"))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&changePasswordPayload)
	if err != nil {
		log.Println("Error while decoding the request body:", err)
		utils.RequestErrorHandler(w, err)
		return
	}

	if changePasswordPayload.Password == "" {
		utils.RequestErrorHandler(w, errors.New("password required, for chaning password common it's common sense"))
		return
	}

	//* Normalize the password
	normalizedPassword := strings.TrimSpace(changePasswordPayload.Password)

	//* Hash the password
	hashedPassword, err := utils.HashPassword(normalizedPassword)

	if err != nil {
		log.Println("Error while hashing password", err)
		utils.InternalErrorHandler(w)
		return
	}

	userRespository := models.UserRepository{DB: db.DB}

	err = userRespository.UpdateUserPassword(userId, hashedPassword)

	if err != nil {
		log.Println("Error while updating the password", err)
		utils.InternalErrorHandler(w)
		return
	}

	utils.OkResponseHandler(w, nil, "Changed password successfully!")
	return

}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var loginResponse LoginResponse

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.RequestErrorHandler(w, err)
		return
	}

	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	user.Password = strings.TrimSpace(user.Password)

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
	log.Println("Password is => ", isPasswordValid)
	if !isPasswordValid {
		utils.NotFoundErrorHandler(w, errors.New("email or password incorrect!"))
		return
	}

	//* Update the last login information
	err = userRepository.UpdateLastLoginAt(userDetails.ID)
	if err != nil {
		log.Println(err)
		//* Not as important field as others
	}

	jwtToken, err := utils.GenerateJwt(userDetails.ID)

	if err != nil {
		log.Println(err)
		utils.InternalErrorHandler(w)
		return
	}

	loginResponse.AccessToken = jwtToken

	utils.OkResponseHandler(w, loginResponse, "User logged in successfully")
	return
}
