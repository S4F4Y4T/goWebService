package service

import (
	"context"

	"github.com/S4F4Y4T/goWebService/internal/model"
)

type UserService struct {
	repo model.UserRepository
}

func NewUserService(repo model.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	return s.repo.Create(ctx, req)
}

func (s *UserService) Update(ctx context.Context, req *model.UpdateUserRequest) (*model.User, error) {
	return s.repo.Update(ctx, req)
}

func (s *UserService) Delete(ctx context.Context, req *model.DeleteUserRequest) error {
	return s.repo.Delete(ctx, req)
}

func (s *UserService) FindByID(ctx context.Context, req *model.GetUserRequest) (*model.User, error) {
	return s.repo.FindByID(ctx, req)
}

func (s *UserService) FindAll(ctx context.Context, req *model.GetUsersRequest) (*model.GetUsersResponse, error) {
	return s.repo.FindAll(ctx, req)
}
