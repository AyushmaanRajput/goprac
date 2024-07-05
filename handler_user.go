package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (app *App) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Age       uint8  `json:"age"`
		Password  string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	fmt.Println("FirstName:", params.FirstName)
	fmt.Println("LastName:", params.LastName)
	fmt.Println("Email:", params.Email)
	fmt.Println("Age:", params.Age)
	fmt.Println("Password:", params.Password)
	if params.FirstName == "" || params.LastName == "" || params.Email == "" || params.Password == "" {
		RespondWithError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	var existingUser User
	if err := app.DB.Where("email = ?", params.Email).First(&existingUser).Error; err == nil {
		RespondWithError(w, http.StatusConflict, "User with the same email already exists")
		return
	} else if err != gorm.ErrRecordNotFound {
		RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	user := User{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Age:       params.Age,
		Password:  string(hashedPassword),
	}

	result := app.DB.Create(&user)
	if result.Error != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	RespondWithJSON(w, http.StatusCreated, CreateUserResponseJSON(user, "User Created Successfully"))
}

func (app *App) handleGetUserById(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")

	var existingUser User
	if err := app.DB.First(&existingUser, "id = ?", userId).Error; err != nil {
		RespondWithError(w, http.StatusNotFound, fmt.Sprintf("User doesn't exist with id : %s", userId))
		return
	}

	RespondWithJSON(w, http.StatusOK, CreateUserResponseJSON(existingUser, "User Found Successfully"))
}
