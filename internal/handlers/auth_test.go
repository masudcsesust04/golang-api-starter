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
	mock "github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	mockDB := new(mocks.MockDB)
	handler := NewUserHandler(nil)
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
	handler := NewUserHandler(nil)
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