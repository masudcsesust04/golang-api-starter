package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/masudcsesust04/golang-jwt-auth/internal/mocks"
	"github.com/masudcsesust04/golang-jwt-auth/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetUsers(t *testing.T) {
	mockDB := new(mocks.MockDB)
	handler := NewUserHandler(nil)
	handler.dbImpl = mockDB

	users := []*models.User{
		{ID: 1, FirstName: "User1", LastName: "Test", Email: "user1@example.com"},
		{ID: 2, FirstName: "User2", LastName: "Test", Email: "user2@example.com"},
	}

	mockDB.On("GetAllUsers").Return(users, nil)

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	handler.GetUsers(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var returnedUsers []*models.User
	json.Unmarshal(w.Body.Bytes(), &returnedUsers)

	assert.Equal(t, users, returnedUsers)
	mockDB.AssertExpectations(t)
}

func TestGetUser(t *testing.T) {
	mockDB := new(mocks.MockDB)
	handler := NewUserHandler(nil)
	handler.dbImpl = mockDB

	user := &models.User{ID: 1, FirstName: "User1", LastName: "Test", Email: "user1@example.com"}

	mockDB.On("GetUserByID", int64(1)).Return(user, nil)

	req := httptest.NewRequest("GET", "/users/1", nil)
	w := httptest.NewRecorder()

	// Need to use mux router to handle path variables
	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", handler.GetUser)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var returnedUser models.User
	json.Unmarshal(w.Body.Bytes(), &returnedUser)

	assert.Equal(t, *user, returnedUser)
	mockDB.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	mockDB := new(mocks.MockDB)
	handler := NewUserHandler(nil)
	handler.dbImpl = mockDB

	user := &models.User{ID: 1, FirstName: "Updated", LastName: "User", Email: "updated@example.com"}

	mockDB.On("UpdateUser", user).Return(nil)

	jsonBody, _ := json.Marshal(user)
	req := httptest.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", handler.UpdateUser)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockDB.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockDB := new(mocks.MockDB)
	handler := NewUserHandler(nil)
	handler.dbImpl = mockDB

	mockDB.On("DeleteUser", int64(1)).Return(nil)

	req := httptest.NewRequest("DELETE", "/users/1", nil)
	w := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", handler.DeleteUser)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockDB.AssertExpectations(t)
}
