package repository

import (
	"errors"
	"sync"

	"github.com/S4F4Y4T/goWebService/internal/model"
)

type userRepository struct {
	mu     sync.RWMutex
	users  map[int]model.User
	nextID int
}

func NewUserRepository() model.UserRepository {
	return &userRepository{
		users:  make(map[int]model.User),
		nextID: 1,
	}
}

func (r *userRepository) Create(req *model.CreateUserRequest) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user := model.User{
		ID:    r.nextID,
		Name:  req.Name,
		Email: req.Email,
	}
	r.users[r.nextID] = user
	r.nextID++

	return &user, nil
}

func (r *userRepository) Update(req *model.UpdateUserRequest) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[req.ID]
	if !ok {
		return nil, errors.New("user not found")
	}

	user.Name = req.Name
	r.users[req.ID] = user

	return &user, nil
}

func (r *userRepository) Delete(req *model.DeleteUserRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[req.ID]; !ok {
		return errors.New("user not found")
	}

	delete(r.users, req.ID)
	return nil
}

func (r *userRepository) FindByID(req *model.GetUserRequest) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[req.ID]
	if !ok {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (r *userRepository) FindAll(req *model.GetUsersRequest) (*model.GetUsersResponse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]model.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return &model.GetUsersResponse{
		Users:  users,
		Total:  len(users),
		Limit:  req.Limit,
		Offset: req.Offset,
	}, nil
}
