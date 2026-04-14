package user

// CreateUserRequest is the DTO for creating a new user
type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
}

// UpdateUserRequest is the DTO for updating an existing user
type UpdateUserRequest struct {
	ID    uint   `json:"id"`
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
}

// GetUsersResponse is the DTO for paginated user list responses
type GetUsersResponse struct {
	Users  []User `json:"users"`
	Total  int64  `json:"total"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}
