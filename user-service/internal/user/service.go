package user

import (
	"context"

	"github.com/S4F4Y4T/goWebService/pkg/apperror"
)

// Service implements user application logic.
type Service struct {
	repo UserRepository
}

func NewService(repo UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req *CreateUserRequest) (*User, error) {
	u, err := NewUser(req.Name, req.Email)
	if err != nil {
		return nil, apperror.New(apperror.BadRequest, err.Error(), err)
	}

	existing, err := s.repo.FindByEmail(ctx, string(u.Email))
	if err != nil {
		return nil, apperror.New(apperror.Internal, "failed to check email uniqueness", err)
	}
	if existing != nil {
		return nil, apperror.New(apperror.BadRequest, "email already taken", nil)
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, apperror.New(apperror.Internal, "failed to create user", err)
	}
	return u, nil
}

func (s *Service) Update(ctx context.Context, req *UpdateUserRequest) (*User, error) {
	u, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, apperror.New(apperror.NotFound, "user not found", err)
	}
	if err := u.UpdateName(req.Name); err != nil {
		return nil, apperror.New(apperror.BadRequest, err.Error(), err)
	}
	if err := u.UpdateEmail(req.Email); err != nil {
		return nil, apperror.New(apperror.BadRequest, err.Error(), err)
	}
	if err := s.repo.Update(ctx, u); err != nil {
		return nil, apperror.New(apperror.Internal, "failed to update user", err)
	}
	return u, nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return apperror.New(apperror.Internal, "failed to delete user", err)
	}
	return nil
}

func (s *Service) FindByID(ctx context.Context, id uint) (*User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.New(apperror.NotFound, "user not found", err)
	}
	return u, nil
}

func (s *Service) FindAll(ctx context.Context, limit, offset int) (*GetUsersResponse, error) {
	users, total, err := s.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, apperror.New(apperror.Internal, "failed to list users", err)
	}
	return &GetUsersResponse{Users: users, Total: total, Limit: limit, Offset: offset}, nil
}
