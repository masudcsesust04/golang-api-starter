package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/masudcsesust04/golang-jwt-auth/internal/models"
	"github.com/masudcsesust04/golang-jwt-auth/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthDBInterface interface {
	RegisterUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	CreateRefreshToken(refreshToken *models.RefreshToken) error
	GetRefreshToken(userID int64) (*models.RefreshToken, error)
	DeleteRefreshToken(userID int64) error
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the login response payload
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var jwtKey = []byte(os.Getenv("JWT_SECRET")) // Replace with your secret key

// Claims represents the JWT claims
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthHandler struct {
	dbImpl AuthDBInterface
}

func NewAuthHandler(user *models.User) *AuthHandler {
	return &AuthHandler{dbImpl: user}
}

// Register handles POST /auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusInternalServerError)
		return
	}

	// Validate the user struct
	if validationErrors := utils.ValidateStruct(user); validationErrors != nil {
		http.Error(w, "Validation failed: "+strings.Join(validationErrors, ", "), http.StatusBadRequest)
		return
	}

	err := h.dbImpl.RegisterUser(&user)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully.",
	})
}

// Login handles POST /auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := h.dbImpl.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	rawSecureToken, err := utils.GenerateSecureToken(32)
	if err != nil {
		http.Error(w, "Secure token generation error", http.StatusInternalServerError)
		return
	}

	hashSecureToken, _ := utils.HashToken(rawSecureToken)
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     hashSecureToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}

	err = h.dbImpl.CreateRefreshToken(refreshToken)
	if err != nil {
		fmt.Printf("Error creating refresh token: %v\n", err)
		http.Error(w, "Failed to create refresh token", http.StatusInternalServerError)
		return
	}

	resp := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		AccessToken:  accessToken,
		RefreshToken: rawSecureToken,
	}

	json.NewEncoder(w).Encode(resp)
}

type RefreshRequest struct {
	UserID       int64  `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshToken handles POST /token/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	refreshToken, err := h.dbImpl.GetRefreshToken(req.UserID)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	err = utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		http.Error(w, "Expired refresh token", http.StatusUnauthorized)
		return
	}

	if utils.CompareToken(refreshToken.Token, req.RefreshToken) != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	accessToken, err := utils.GenerateAccessToken(req.UserID)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	// TODO: Optionally rotate refresh token here and update expire_at
	// response with new refresh token
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
}

// Logout handles POST /auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// For logout, typically the client deletes the tokens.
	// Here, we expect the user_id in the request body to delete refresh token from DB.

	type LogoutRequest struct {
		UserID int64 `json:"user_id"`
	}

	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.dbImpl.DeleteRefreshToken(req.UserID)
	if err != nil {
		http.Error(w, "Failed to logout: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
