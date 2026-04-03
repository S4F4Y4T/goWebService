package service

import (
	"github.com/S4F4Y4T/goWebService/internal/model"
)

type UserService struct {
	repo model.UserRepository
}

func NewUserService(repo model.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(req *model.CreateUserRequest) (*model.User, error) {
	return s.repo.Create(req)
}

func (s *UserService) Update(req *model.UpdateUserRequest) (*model.User, error) {
	return s.repo.Update(req)
}

func (s *UserService) Delete(req *model.DeleteUserRequest) error {
	return s.repo.Delete(req)
}

func (s *UserService) FindByID(req *model.GetUserRequest) (*model.User, error) {
	return s.repo.FindByID(req)
}

func (s *UserService) FindAll(req *model.GetUsersRequest) (*model.GetUsersResponse, error) {
	return s.repo.FindAll(req)
}
