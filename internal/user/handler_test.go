package user_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/S4F4Y4T/goWebService/internal/user"
)

type MockUserService struct {
	CreateFn   func(ctx context.Context, req *user.CreateUserRequest) (*user.User, error)
	UpdateFn   func(ctx context.Context, req *user.UpdateUserRequest) (*user.User, error)
	DeleteFn   func(ctx context.Context, id uint) error
	FindByIDFn func(ctx context.Context, id uint) (*user.User, error)
	FindAllFn  func(ctx context.Context, limit, offset int) (*user.GetUsersResponse, error)
}

func (m *MockUserService) Create(ctx context.Context, req *user.CreateUserRequest) (*user.User, error) {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, req)
	}
	return nil, nil
}

func (m *MockUserService) Update(ctx context.Context, req *user.UpdateUserRequest) (*user.User, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, req)
	}
	return nil, nil
}

func (m *MockUserService) Delete(ctx context.Context, id uint) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *MockUserService) FindByID(ctx context.Context, id uint) (*user.User, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockUserService) FindAll(ctx context.Context, limit, offset int) (*user.GetUsersResponse, error) {
	if m.FindAllFn != nil {
		return m.FindAllFn(ctx, limit, offset)
	}
	return nil, nil
}

func TestHandler_CreateUser(t *testing.T) {
	mockService := &MockUserService{
		CreateFn: func(ctx context.Context, req *user.CreateUserRequest) (*user.User, error) {
			u, _ := user.NewUser(req.Name, req.Email)
			u.ID = 1
			return u, nil
		},
	}
	handler := user.NewHandler(mockService)

	reqBody := `{"name":"John Doe","email":"john@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateUser(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, res.StatusCode)
	}
}

func TestHandler_GetUser(t *testing.T) {
	mockService := &MockUserService{
		FindByIDFn: func(ctx context.Context, id uint) (*user.User, error) {
			u, _ := user.NewUser("John Doe", "john@example.com")
			u.ID = id
			return u, nil
		},
	}
	handler := user.NewHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.GetUser(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.StatusCode)
	}
}

func TestHandler_GetUsers(t *testing.T) {
	mockService := &MockUserService{
		FindAllFn: func(ctx context.Context, limit, offset int) (*user.GetUsersResponse, error) {
			u, _ := user.NewUser("John Doe", "john@example.com")
			u.ID = 1
			return &user.GetUsersResponse{
				Users: []user.User{*u},
				Total: 1,
			}, nil
		},
	}
	handler := user.NewHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	handler.GetUsers(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.StatusCode)
	}
}

func TestHandler_DeleteUser(t *testing.T) {
	mockService := &MockUserService{
		DeleteFn: func(ctx context.Context, id uint) error {
			return nil
		},
	}
	handler := user.NewHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.DeleteUser(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.StatusCode)
	}
}
