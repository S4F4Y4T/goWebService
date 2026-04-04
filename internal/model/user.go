package model

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
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
	Create(req *CreateUserRequest) (*User, error)
	Update(req *UpdateUserRequest) (*User, error)
	Delete(req *DeleteUserRequest) error
	FindByID(req *GetUserRequest) (*User, error)
	FindAll(req *GetUsersRequest) (*GetUsersResponse, error)
}
