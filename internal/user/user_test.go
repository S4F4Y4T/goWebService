package user_test

import (
	"strings"
	"testing"

	"github.com/S4F4Y4T/goWebService/internal/user"
)

func TestNewUser_Valid(t *testing.T) {
	u, err := user.NewUser("John Doe", "john@example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if u.Name != "John Doe" {
		t.Errorf("expected Name 'John Doe', got %v", u.Name)
	}
	if string(u.Email) != "john@example.com" {
		t.Errorf("expected Email 'john@example.com', got %v", u.Email)
	}
}

func TestNewUser_InvalidName(t *testing.T) {
	_, err := user.NewUser("J", "john@example.com")
	if err == nil {
		t.Fatal("expected error for too short name, got nil")
	}
	if err != user.ErrInvalidName {
		t.Errorf("expected ErrInvalidName, got %v", err)
	}

	longName := strings.Repeat("A", 101)
	_, err = user.NewUser(longName, "john@example.com")
	if err == nil {
		t.Fatal("expected error for too long name, got nil")
	}
	if err != user.ErrInvalidName {
		t.Errorf("expected ErrInvalidName, got %v", err)
	}
}

func TestNewUser_InvalidEmail(t *testing.T) {
	_, err := user.NewUser("John Doe", "invalid-email")
	if err == nil {
		t.Fatal("expected error for invalid email, got nil")
	}
	if err != user.ErrInvalidEmail {
		t.Errorf("expected ErrInvalidEmail, got %v", err)
	}
}
