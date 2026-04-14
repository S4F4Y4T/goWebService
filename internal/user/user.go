package user

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/S4F4Y4T/goWebService/internal/shared/domain"
)

var (
	ErrInvalidEmail = errors.New("invalid email format")
	ErrInvalidName  = errors.New("name must be between 2 and 100 characters")
)

// Email is a Value Object representing a validated email address
type Email string

func NewEmail(v string) (Email, error) {
	v = strings.TrimSpace(strings.ToLower(v))
	// Simple email regex for demonstration
	re := regexp.MustCompile(`^[a-z0-0._%+-]+@[a-z0-0.-]+\.[a-z]{2,}$`)
	if !re.MatchString(v) {
		return "", ErrInvalidEmail
	}
	return Email(v), nil
}

func (e Email) String() string {
	return string(e)
}

// User is the Aggregate Root for the User domain
type User struct {
	domain.AggregateRoot `json:"-" gorm:"-"`
	ID                   uint      `json:"id"`
	Name                 string    `json:"name"`
	Email                Email     `json:"email"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// NewUser is a factory function for creating a valid User entity
func NewUser(name, emailStr string) (*User, error) {
	u := &User{}
	if err := u.UpdateName(name); err != nil {
		return nil, err
	}
	if err := u.UpdateEmail(emailStr); err != nil {
		return nil, err
	}
	return u, nil
}

// Domain Methods (Rich Domain Model)

func (u *User) UpdateName(name string) error {
	name = strings.TrimSpace(name)
	if len(name) < 2 || len(name) > 100 {
		return ErrInvalidName
	}
	u.Name = name
	return nil
}

func (u *User) UpdateEmail(emailStr string) error {
	email, err := NewEmail(emailStr)
	if err != nil {
		return err
	}
	u.Email = email
	return nil
}

// UserRepository defines the persistence contract for Users
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*User, error)
	FindAll(ctx context.Context, limit, offset int) ([]User, int64, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}

// API DTOs
type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserRequest struct {
	ID    uint   `json:"id"`
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
}

type GetUsersResponse struct {
	Users  []User `json:"users"`
	Total  int64  `json:"total"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}
