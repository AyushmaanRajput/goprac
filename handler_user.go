package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func (app *App) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		FirstName string
		LastName  string
		Email     string
		Age       uint8
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	var existingUser User
	if err := app.DB.Where("email = ?", params.Email).First(&existingUser).Error; err == nil {
		respondWithError(w, http.StatusConflict, "User with the same email already exists")
		return
	} else if err != gorm.ErrRecordNotFound {
		respondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	user := User{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Age:       params.Age,
	}

	result := app.DB.Create(&user)
	if result.Error != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	respondWithJSON(w, http.StatusCreated, createUserResponseJSON(user, "User Created Successfully"))
}
