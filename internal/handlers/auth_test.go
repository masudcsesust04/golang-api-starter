package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/masudcsesust04/golang-jwt-auth/internal/mocks"
	"github.com/masudcsesust04/golang-jwt-auth/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	mockDB := new(mocks.MockDB)
	handler := NewAuthHandler(nil)
	handler.dbImpl = mockDB

	user := &models.User{FirstName: "New", LastName: "User", Email: "new@example.com", Password: "password123", Status: "active", PhoneNumber: "+1234567890"}

	mockDB.On("RegisterUser", user).Return(nil)

	jsonBody, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()

	handler.Register(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockDB.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	mockDB := new(mocks.MockDB)
	handler := NewAuthHandler(nil)
	handler.dbImpl = mockDB

	// Mock user data
	password := "testpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &models.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
	}

	loginReq := LoginRequest{Email: "test@example.com", Password: password}
	mockDB.On("GetUserByEmail", loginReq.Email).Return(user, nil).Once()
	mockDB.On("CreateRefreshToken", mock.AnythingOfType("*models.RefreshToken")).Return(nil).Once()

	jsonBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()

	handler.Login(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var loginResp LoginResponse
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	assert.NotEmpty(t, loginResp.AccessToken)
	assert.NotEmpty(t, loginResp.RefreshToken)
	mockDB.AssertExpectations(t)
}

func TestRefreshToken(t *testing.T) {
	mockDB := new(mocks.MockDB)
	handler := NewAuthHandler(nil)
	handler.dbImpl = mockDB

	// Mock refresh token data
	rawToken := "raw_refresh_token"
	hashedToken, _ := bcrypt.GenerateFromPassword([]byte(rawToken), bcrypt.DefaultCost)
	refreshToken := &models.RefreshToken{
		UserID:    1,
		Token:     string(hashedToken),
		ExpiresAt: time.Now().Add(time.Hour),
	}

	refreshReq := RefreshRequest{UserID: 1, RefreshToken: rawToken}
	mockDB.On("GetRefreshToken", refreshReq.UserID).Return(refreshToken, nil).Once()

	jsonBody, _ := json.Marshal(refreshReq)
	req := httptest.NewRequest("POST", "/token/refresh", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()

	handler.RefreshToken(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp["access_token"])
	mockDB.AssertExpectations(t)
}
