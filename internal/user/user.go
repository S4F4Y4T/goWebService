package user

import (
	"context"
	"time"
)

// User is the Aggregate Root for the User domain
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	Email string `json:"email" validate:"required,email,unique_email"`
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
