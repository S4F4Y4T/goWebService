package model

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DeleteUserRequest struct {
	ID int `json:"id"`
}

type GetUserRequest struct {
	ID int `json:"id"`
}

type GetUsersRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type GetUsersResponse struct {
	Users  []User `json:"users"`
	Total  int    `json:"total"`
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
