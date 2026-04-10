package service

import (
	"context"
	"errors"

	"github.com/S4F4Y4T/goWebService/internal/model"
	"github.com/S4F4Y4T/goWebService/pkg/apperror"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

var tracer = otel.Tracer("user-service")

type UserService struct {
	repo model.UserRepository
}

func NewUserService(repo model.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	ctx, span := tracer.Start(ctx, "service.create_user")
	defer span.End()

	user, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, apperror.New(apperror.Internal, "failed to create user", err)
	}
	return user, nil
}

func (s *UserService) Update(ctx context.Context, req *model.UpdateUserRequest) (*model.User, error) {
	user, err := s.repo.Update(ctx, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "user not found" {
			return nil, apperror.New(apperror.NotFound, "user not found", err)
		}
		return nil, apperror.New(apperror.Internal, "failed to update user", err)
	}
	return user, nil
}

func (s *UserService) Delete(ctx context.Context, req *model.DeleteUserRequest) error {
	err := s.repo.Delete(ctx, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "user not found" {
			return apperror.New(apperror.NotFound, "user not found", err)
		}
		return apperror.New(apperror.Internal, "failed to delete user", err)
	}
	return nil
}

func (s *UserService) FindByID(ctx context.Context, req *model.GetUserRequest) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "user not found" {
			return nil, apperror.New(apperror.NotFound, "user not found", err)
		}
		return nil, apperror.New(apperror.Internal, "failed to find user", err)
	}
	return user, nil
}

func (s *UserService) FindAll(ctx context.Context, req *model.GetUsersRequest) (*model.GetUsersResponse, error) {
	res, err := s.repo.FindAll(ctx, req)
	if err != nil {
		return nil, apperror.New(apperror.Internal, "failed to list users", err)
	}
	return res, nil
}

func (s *UserService) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, apperror.New(apperror.Internal, "failed to find user by email", err)
	}
	return user, nil
}
