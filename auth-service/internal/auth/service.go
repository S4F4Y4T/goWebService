package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/S4F4Y4T/goWebService/pkg/apperror"
	"github.com/S4F4Y4T/goWebService/pkg/jwtutil"
	"golang.org/x/crypto/bcrypt"
)

// Service handles authentication business logic.
type Service struct {
	repo           CredentialRepository
	userServiceURL string
	httpClient     *http.Client
}

func NewService(repo CredentialRepository, userServiceURL string) *Service {
	return &Service{
		repo:           repo,
		userServiceURL: userServiceURL,
		httpClient:     &http.Client{},
	}
}

// Register creates a user profile (via user-service) and stores the hashed password.
func (s *Service) Register(req *RegisterRequest) (*AuthResponse, error) {
	// 1. Check if credentials already exist for this email.
	existing, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, apperror.New(apperror.Internal, "failed to check existing credentials", err)
	}
	if existing != nil {
		return nil, apperror.New(apperror.BadRequest, "email already registered", nil)
	}

	// 2. Create user profile in user-service.
	userID, err := s.createUserProfile(req.Name, req.Email)
	if err != nil {
		return nil, err
	}

	// 3. Hash the password.
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.New(apperror.Internal, "failed to hash password", err)
	}

	// 4. Persist credentials.
	cred := &Credential{
		UserID:       userID,
		Email:        req.Email,
		PasswordHash: string(hash),
	}
	if err := s.repo.Save(cred); err != nil {
		return nil, apperror.New(apperror.Internal, "failed to save credentials", err)
	}

	// 5. Issue JWT.
	token, err := jwtutil.GenerateToken(userID, req.Email)
	if err != nil {
		return nil, apperror.New(apperror.Internal, "failed to generate token", err)
	}

	return &AuthResponse{
		Token: token,
		User:  map[string]any{"id": userID, "email": req.Email, "name": req.Name},
	}, nil
}

// Login verifies credentials and returns a JWT.
func (s *Service) Login(req *LoginRequest) (*AuthResponse, error) {
	cred, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, apperror.New(apperror.Internal, "failed to find credentials", err)
	}
	if cred == nil {
		return nil, apperror.New(apperror.Unauthorized, "invalid email or password", nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(cred.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperror.New(apperror.Unauthorized, "invalid email or password", nil)
	}

	token, err := jwtutil.GenerateToken(cred.UserID, cred.Email)
	if err != nil {
		return nil, apperror.New(apperror.Internal, "failed to generate token", err)
	}

	return &AuthResponse{
		Token: token,
		User:  map[string]any{"id": cred.UserID, "email": cred.Email},
	}, nil
}

// createUserProfile calls user-service to create the user record.
func (s *Service) createUserProfile(name, email string) (uint, error) {
	body, _ := json.Marshal(map[string]string{"name": name, "email": email})
	resp, err := s.httpClient.Post(
		fmt.Sprintf("%s/users/", s.userServiceURL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return 0, apperror.New(apperror.Internal, "could not reach user-service", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return 0, apperror.New(apperror.Internal,
			fmt.Sprintf("user-service returned %d: %s", resp.StatusCode, string(raw)), nil)
	}

	// Parse the user-service response: { "success": true, "data": { "id": 1, ... } }
	var envelope struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return 0, apperror.New(apperror.Internal, "failed to parse user-service response", err)
	}
	if envelope.Data.ID == 0 {
		return 0, apperror.New(apperror.Internal, "user-service returned unexpected response", nil)
	}
	return envelope.Data.ID, nil
}
