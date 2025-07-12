package models

import (
	"os"
	"testing"
	"time"

	"github.com/masudcsesust04/golang-jwt-auth/internal/config"
)

func TestMain(m *testing.M) {
	cleanup := config.SetupTestDB(m)
	defer cleanup()

	os.Exit(m.Run())
}

func TestCreateAndGetUser(t *testing.T) {
	user := &User{
		FirstName:   "Test",
		LastName:    "User",
		PhoneNumber: "1234567890",
		Email:       "testuser@example.com",
		Status:      "active",
		Password:    "password123",
	}

	err := user.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// Use GetUserByEmail instead of GetUserByUsername
	gotUser, err := user.GetUserByEmail(user.Email)
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}
	if gotUser == nil || gotUser.Email != user.Email {
		t.Fatalf("GetUserByEmail returned wrong user")
	}
}

func TestUpdateUser(t *testing.T) {
	user := &User{
		FirstName:   "Update",
		LastName:    "User",
		PhoneNumber: "0987654321",
		Email:       "updateuser@example.com",
		Status:      "active",
		Password:    "password123",
	}

	err := user.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	user.Email = "newemail@example.com"
	err = user.UpdateUser(user)
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}

	updatedUser, err := user.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if updatedUser.Email != "newemail@example.com" {
		t.Fatalf("UpdateUser did not update email")
	}
}

func TestDeleteUser(t *testing.T) {
	user := &User{
		FirstName:   "Delete",
		LastName:    "User",
		PhoneNumber: "1112223333",
		Email:       "deleteuser@example.com",
		Status:      "active",
		Password:    "password123",
	}

	err := user.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	err = user.DeleteUser(user.ID)
	if err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}

	deletedUser, err := user.GetUserByID(user.ID)
	if err == nil && deletedUser != nil {
		t.Fatalf("DeleteUser did not delete user")
	}
}

func TestCreateAndDeleteRefreshToken(t *testing.T) {
	user := &User{
		FirstName:   "Token",
		LastName:    "User",
		PhoneNumber: "4445556666",
		Email:       "tokenuser@example.com",
		Password:    "password123",
	}

	err := user.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	rt := &RefreshToken{
		UserID:    user.ID,
		Token:     "testtoken",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	err = user.CreateRefreshToken(rt)
	if err != nil {
		t.Fatalf("CreateRefreshToken failed: %v", err)
	}

	gotRT, err := user.GetRefreshToken(rt.UserID)
	if err != nil {
		t.Fatalf("GetRefreshToken failed: %v", err)
	}
	if gotRT == nil || gotRT.Token != rt.Token {
		t.Fatalf("GetRefreshToken returned wrong token")
	}

	err = user.DeleteRefreshToken(rt.UserID)
	if err != nil {
		t.Fatalf("DeleteRefreshToken failed: %v", err)
	}

	deletedRT, err := user.GetRefreshToken(rt.UserID)
	if err == nil && deletedRT != nil {
		t.Fatalf("DeleteRefreshToken did not delete token")
	}
}

func TestGetAllUsers(t *testing.T) {
	user := &User{}
	users, err := user.GetAllUsers()
	if err != nil {
		t.Fatalf("GetAllUsers failed: %v", err)
	}
	if len(users) == 0 {
		t.Fatalf("GetAllUsers returned empty list")
	}
}