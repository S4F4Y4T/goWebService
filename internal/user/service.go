package user

import (
	"context"

	"github.com/S4F4Y4T/goWebService/pkg/apperror"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("user-service")

type Service struct {
	repo UserRepository
}

func NewService(repo UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req *CreateUserRequest) (*User, error) {
	ctx, span := tracer.Start(ctx, "user.service.create")
	defer span.End()

	u := &User{
		Name:  req.Name,
		Email: req.Email,
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

	u.Name = req.Name
	u.Email = req.Email

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

	return &GetUsersResponse{
		Users:  users,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *Service) FindByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.FindByEmail(ctx, email)
}
