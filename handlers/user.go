package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/daluisgarcia/golang-rest-websockets/models"
	"github.com/daluisgarcia/golang-rest-websockets/repositories"
	"github.com/daluisgarcia/golang-rest-websockets/server"
	"github.com/segmentio/ksuid"
)

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = SignUpRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userId, err := ksuid.NewRandom()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user := models.User{
			Id:       userId.String(),
			Email:    request.Email,
			Password: request.Password,
		}

		err = repositories.InsertUser(r.Context(), &user)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(SignUpResponse{
			Id:    user.Id,
			Email: user.Email,
		})
	}
}
