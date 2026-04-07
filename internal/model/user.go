package model

import (
	"context"
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email,unique_email"`
}

type UpdateUserRequest struct {
	ID    uint   `json:"id"`
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
}

type DeleteUserRequest struct {
	ID uint `json:"id"`
}

type GetUserRequest struct {
	ID uint `json:"id"`
}

type GetUsersRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type GetUsersResponse struct {
	Users  []User `json:"users"`
	Total  int64  `json:"total"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type UserRepository interface {
	Create(ctx context.Context, req *CreateUserRequest) (*User, error)
	Update(ctx context.Context, req *UpdateUserRequest) (*User, error)
	Delete(ctx context.Context, req *DeleteUserRequest) error
	FindByID(ctx context.Context, req *GetUserRequest) (*User, error)
	FindAll(ctx context.Context, req *GetUsersRequest) (*GetUsersResponse, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}
