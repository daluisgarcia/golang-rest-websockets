package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/daluisgarcia/golang-rest-websockets/middleware"
	"github.com/daluisgarcia/golang-rest-websockets/models"
	"github.com/daluisgarcia/golang-rest-websockets/repositories"
	"github.com/daluisgarcia/golang-rest-websockets/server"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
)

type PostRequest struct {
	PostContent string `json:"postContent"`
}

type PostResponse struct {
	Id          string `json:"id"`
	PostContent string `json:"postContent"`
}

func InsertPostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request PostRequest
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

			// Build a message to be sent to the websocket
			var postWebSocketMessage = models.WebSocketMessage{
				Type:    "Post Created",
				Payload: post,
			}

			// Notifies through websockets that a new post has been created
			s.Hub().Broadcast(postWebSocketMessage, nil)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(PostResponse{
				Id:          post.Id,
				PostContent: post.PostContent,
			})
		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

	}
}

func GetPostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		post, err := repositories.FindPostById(r.Context(), params["id"])

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if post == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(post)
	}
}

func UpdatePostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := middleware.GetJwtTokenFromHeader(s, r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if _, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			params := mux.Vars(r)
			var request PostRequest
			err := json.NewDecoder(r.Body).Decode(&request)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			post, err := repositories.FindPostById(r.Context(), params["id"])

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if post == nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			post.PostContent = request.PostContent

			err = repositories.UpdatePost(r.Context(), post)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(PostResponse{
				Id:          post.Id,
				PostContent: post.PostContent,
			})
		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

	}
}

func DeletePostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := middleware.GetJwtTokenFromHeader(s, r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			params := mux.Vars(r)
			post, err := repositories.FindPostById(r.Context(), params["id"])

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if post == nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			err = repositories.DeletePost(r.Context(), post.Id, claims.UserId)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

	}
}

func ListPostsHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := middleware.GetJwtTokenFromHeader(s, r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			pageStr := r.URL.Query().Get("page")

			var page = uint64(0)

			if pageStr != "" {
				page, err = strconv.ParseUint(pageStr, 10, 64)

				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}

			posts, err := repositories.ListPosts(r.Context(), page, claims.UserId)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(posts)
		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
	}
}
