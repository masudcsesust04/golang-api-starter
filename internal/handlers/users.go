package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/masudcsesust04/golang-jwt-auth/internal/models"
)

type UserDBInterface interface {
	GetAllUsers() ([]*models.User, error)
	GetUserByID(id int64) (*models.User, error)
	GetUserByEmail(emaio string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(id int64) error
	CreateRefreshToken(refreshToken *models.RefreshToken) error
	GetRefreshToken(userID int64) (*models.RefreshToken, error)
	DeleteRefreshToken(userID int64) error
}

type UserHandler struct {
	dbImpl UserDBInterface
}

func NewUserHandler(user *models.User) *UserHandler {
	return &UserHandler{dbImpl: user}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.dbImpl.GetAllUsers()
	if err != nil {
		http.Error(w, "failed to get users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) CreateUsers(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusInternalServerError)
		return
	}

	err := h.dbImpl.CreateUser(&user)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUser handle GET /users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIdStr := vars["id"]
	id, err := strconv.ParseInt(userIdStr, 10, 64)

	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
	}

	user, err := h.dbImpl.GetUserByID(id)
	if err != nil {
		http.Error(w, "Failed to get user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIdStr := vars["id"]
	id, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user.ID = id
	err = h.dbImpl.UpdateUser(&user)
	if err != nil {
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIdStr := vars["id"]
	id, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
	}

	err = h.dbImpl.DeleteUser(id)
	if err != nil {
		http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
