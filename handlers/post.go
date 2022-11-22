package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/daluisgarcia/golang-rest-websockets/middleware"
	"github.com/daluisgarcia/golang-rest-websockets/models"
	"github.com/daluisgarcia/golang-rest-websockets/repositories"
	"github.com/daluisgarcia/golang-rest-websockets/server"
	"github.com/segmentio/ksuid"
)

type InsertPostRequest struct {
	PostContent string `json:"postContent"`
}

type InserPostResponse struct {
	Id          string `json:"id"`
	PostContent string `json:"postContent"`
}

func InsertPostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request InsertPostRequest
		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token, err := middleware.GetJwtTokenFromHeader(s, r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {

			id, err := ksuid.NewRandom()

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			post := &models.Post{
				Id:          id.String(),
				UserId:      claims.UserId,
				PostContent: request.PostContent,
			}

			err = repositories.InsertPost(r.Context(), post)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(InserPostResponse{
				Id:          post.Id,
				PostContent: post.PostContent,
			})
		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

	}
}
