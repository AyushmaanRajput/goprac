package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var JWTSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func (app *App) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	var existingUser User
	if err := app.DB.Where("email = ?", params.Email).First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			RespondWithError(w, http.StatusNotFound, "User with this email doesn't exist")
		} else {
			RespondWithError(w, http.StatusInternalServerError, "Database error")
		}
		return
	}

	// Validate the password
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(params.Password)); err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": existingUser.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expires after 72 hours
	})

	tokenString, err := token.SignedString(JWTSecretKey)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	response := map[string]interface{}{
		"message": "User Logged In Successfully",
		"token":   tokenString,
		"user":    existingUser,
	}

	RespondWithJSON(w, http.StatusOK, response)
}
