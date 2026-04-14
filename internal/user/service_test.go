package user_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/S4F4Y4T/goWebService/internal/shared/domain"
	"github.com/S4F4Y4T/goWebService/internal/shared/event"
	"github.com/S4F4Y4T/goWebService/internal/user"
)

// MockUserRepository implements user.UserRepository for testing
type MockUserRepository struct {
	CreateFn      func(ctx context.Context, u *user.User) error
	UpdateFn      func(ctx context.Context, u *user.User) error
	DeleteFn      func(ctx context.Context, id uint) error
	FindByIDFn    func(ctx context.Context, id uint) (*user.User, error)
	FindAllFn     func(ctx context.Context, limit, offset int) ([]user.User, int64, error)
	FindByEmailFn func(ctx context.Context, email string) (*user.User, error)
}

func (m *MockUserRepository) Create(ctx context.Context, u *user.User) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, u)
	}
	return nil
}

func (m *MockUserRepository) Update(ctx context.Context, u *user.User) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, u)
	}
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*user.User, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockUserRepository) FindAll(ctx context.Context, limit, offset int) ([]user.User, int64, error) {
	if m.FindAllFn != nil {
		return m.FindAllFn(ctx, limit, offset)
	}
	return nil, 0, nil
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	if m.FindByEmailFn != nil {
		return m.FindByEmailFn(ctx, email)
	}
	return nil, nil
}

func TestService_Create(t *testing.T) {
	mockRepo := &MockUserRepository{
		FindByEmailFn: func(ctx context.Context, email string) (*user.User, error) {
			return nil, nil // Email available
		},
		CreateFn: func(ctx context.Context, u *user.User) error {
			u.ID = 1
			u.CreatedAt = time.Now()
			u.UpdatedAt = time.Now()
			return nil
		},
	}
	dispatcher := event.NewDispatcher()

	var eventFired bool
	dispatcher.Subscribe(user.UserCreatedTopic, func(ctx context.Context, ev domain.DomainEvent) error {
		if _, ok := ev.(user.UserCreated); ok {
			eventFired = true
		}
		return nil
	})

	svc := user.NewService(mockRepo, dispatcher)

	req := &user.CreateUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
	}

	u, err := svc.Create(context.Background(), req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if u == nil {
		t.Fatal("expected user, got nil")
	}
	if u.ID != 1 {
		t.Errorf("expected ID 1, got %v", u.ID)
	}
	if u.Name != "Test User" {
		t.Errorf("expected Name 'Test User', got %v", u.Name)
	}
	if u.Email.String() != "test@example.com" {
		t.Errorf("expected Email 'test@example.com', got %v", u.Email)
	}

	if !eventFired {
		t.Errorf("expected UserCreated event to be dispatched, but it wasn't")
	}
}

func TestService_Create_EmailAlreadyTaken(t *testing.T) {
	mockRepo := &MockUserRepository{
		FindByEmailFn: func(ctx context.Context, email string) (*user.User, error) {
			// Simulate that the email is already in use
			return &user.User{ID: 2, Name: "Existing User"}, nil
		},
	}
	dispatcher := event.NewDispatcher()
	svc := user.NewService(mockRepo, dispatcher)

	req := &user.CreateUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
	}

	u, err := svc.Create(context.Background(), req)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if u != nil {
		t.Fatal("expected nil user, got user")
	}
	if !strings.Contains(err.Error(), "email already taken") {
		t.Errorf("unexpected error message: %v", err.Error())
	}
}

func TestService_Update(t *testing.T) {
	mockRepo := &MockUserRepository{
		FindByIDFn: func(ctx context.Context, id uint) (*user.User, error) {
			if id == 1 {
				u, _ := user.NewUser("Old Name", "old@example.com")
				u.ID = 1
				return u, nil
			}
			return nil, errors.New("not found")
		},
		UpdateFn: func(ctx context.Context, u *user.User) error {
			return nil
		},
	}
	dispatcher := event.NewDispatcher()
	svc := user.NewService(mockRepo, dispatcher)

	req := &user.UpdateUserRequest{
		ID:    1,
		Name:  "New Name",
		Email: "new@example.com",
	}

	u, err := svc.Update(context.Background(), req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if u.Name != "New Name" {
		t.Errorf("expected Name 'New Name', got %v", u.Name)
	}
	if u.Email.String() != "new@example.com" {
		t.Errorf("expected Email 'new@example.com', got %v", u.Email)
	}
}

func TestService_Update_NotFound(t *testing.T) {
	mockRepo := &MockUserRepository{
		FindByIDFn: func(ctx context.Context, id uint) (*user.User, error) {
			return nil, errors.New("not found")
		},
	}
	dispatcher := event.NewDispatcher()
	svc := user.NewService(mockRepo, dispatcher)

	req := &user.UpdateUserRequest{
		ID:    99,
		Name:  "New Name",
		Email: "new@example.com",
	}

	_, err := svc.Update(context.Background(), req)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "user not found") {
		t.Errorf("unexpected error message: %v", err.Error())
	}
}

func TestService_Delete(t *testing.T) {
	mockRepo := &MockUserRepository{
		DeleteFn: func(ctx context.Context, id uint) error {
			return nil
		},
	}
	dispatcher := event.NewDispatcher()
	svc := user.NewService(mockRepo, dispatcher)

	err := svc.Delete(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestService_FindByID(t *testing.T) {
	mockRepo := &MockUserRepository{
		FindByIDFn: func(ctx context.Context, id uint) (*user.User, error) {
			if id == 1 {
				u, _ := user.NewUser("Test Name", "test@example.com")
				u.ID = 1
				return u, nil
			}
			return nil, errors.New("not found")
		},
	}
	dispatcher := event.NewDispatcher()
	svc := user.NewService(mockRepo, dispatcher)

	u, err := svc.FindByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if u.ID != 1 {
		t.Errorf("expected ID 1, got %v", u.ID)
	}

	_, err = svc.FindByID(context.Background(), 99)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestService_FindAll(t *testing.T) {
	mockRepo := &MockUserRepository{
		FindAllFn: func(ctx context.Context, limit, offset int) ([]user.User, int64, error) {
			u, _ := user.NewUser("Test Name", "test@example.com")
			u.ID = 1
			return []user.User{*u}, 1, nil
		},
	}
	dispatcher := event.NewDispatcher()
	svc := user.NewService(mockRepo, dispatcher)

	res, err := svc.FindAll(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Total != 1 {
		t.Errorf("expected Total 1, got %v", res.Total)
	}
	if len(res.Users) != 1 {
		t.Errorf("expected 1 user, got %v", len(res.Users))
	}
}
