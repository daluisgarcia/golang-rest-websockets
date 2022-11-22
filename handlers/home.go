package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/daluisgarcia/golang-rest-websockets/server"
)

type HomeResponse struct {
	Message string `json:"message"` // This strings at the end allows to specify the name of the field in the json response when serializing the struct
	Status  bool   `json:"status"`
}

func HomeHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(HomeResponse{
			Message: "Hello World",
			Status:  true,
		})
	}
}
