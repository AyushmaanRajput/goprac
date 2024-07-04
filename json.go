package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type App struct {
	DB *gorm.DB
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, status int, msg string) {
	if status > 499 {
		log.Println("Responding with 5XX error", msg)
	}

	type errResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, status, errResponse{
		Error: msg,
	})
}

type UserResponse struct {
	Message string `json:"message"`
	User    User   `json:"user"`
}

func createUserResponseJSON(user User, message string) UserResponse {
	return UserResponse{
		Message: message,
		User:    user,
	}
}
