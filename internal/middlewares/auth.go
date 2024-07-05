package middlewares

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/lamichhaneshuvam/todo-pg/internal/utils"
)

func UserAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")

		if bearerToken == "" {
			utils.ForbiddenErrorHandler(w, errors.New("Please login to continue"))
			return
		}

		claims, err := utils.ValidateJWT(bearerToken)

		if err != nil {
			log.Println(err)
			utils.ForbiddenErrorHandler(w, errors.New("Token expired, please login to continue"))
			return
		}

		r.Header.Set("userId", strconv.Itoa(claims.ID))
		next(w, r)
	}
}
