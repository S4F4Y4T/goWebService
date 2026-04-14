package user

import (
	"context"

	"github.com/S4F4Y4T/goWebService/internal/shared/event"
	"github.com/S4F4Y4T/goWebService/pkg/apperror"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("user-service")

type Service struct {
	repo       UserRepository
	dispatcher *event.Dispatcher
}

func NewService(repo UserRepository, dispatcher *event.Dispatcher) *Service {
	return &Service{
		repo:       repo,
		dispatcher: dispatcher,
	}
}

func (s *Service) Create(ctx context.Context, req *CreateUserRequest) (*User, error) {
	ctx, span := tracer.Start(ctx, "user.service.create")
	defer span.End()

	u, err := NewUser(req.Name, req.Email)
	if err != nil {
		return nil, apperror.New(apperror.BadRequest, err.Error(), err)
	}

	// Business Validation: Uniqueness
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

	// Record and Dispatch Event
	u.RecordEvent(NewUserCreated(u.ID, string(u.Email)))
	s.dispatcher.Dispatch(ctx, u.GetEvents())
	u.ClearEvents()

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
