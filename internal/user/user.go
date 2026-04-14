package user

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/S4F4Y4T/goWebService/internal/shared/domain"
)

var (
	ErrInvalidEmail = errors.New("invalid email format")
	ErrInvalidName  = errors.New("name must be between 2 and 100 characters")
)

// Email is a Value Object representing a validated email address
type Email string

func NewEmail(v string) (Email, error) {
	v = strings.TrimSpace(strings.ToLower(v))
	// Fix: [a-z0-9] not [a-z0-0]
	re := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	if !re.MatchString(v) {
		return "", ErrInvalidEmail
	}
	return Email(v), nil
}

func (e Email) String() string {
	return string(e)
}

// User is the Aggregate Root for the User domain.
// AggregateRoot is embedded with json:"-" only — no gorm tags
// because GORM never touches the User entity directly (userSchema handles that).
type User struct {
	domain.AggregateRoot `json:"-"`
	ID                   uint      `json:"id"`
	Name                 string    `json:"name"`
	Email                Email     `json:"email"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// NewUser is a factory function for creating a valid User entity
func NewUser(name, emailStr string) (*User, error) {
	u := &User{}
	if err := u.UpdateName(name); err != nil {
		return nil, err
	}
	if err := u.UpdateEmail(emailStr); err != nil {
		return nil, err
	}
	return u, nil
}

// Domain Methods (Rich Domain Model)

func (u *User) UpdateName(name string) error {
	name = strings.TrimSpace(name)
	if len(name) < 2 || len(name) > 100 {
		return ErrInvalidName
	}
	u.Name = name
	return nil
}

func (u *User) UpdateEmail(emailStr string) error {
	email, err := NewEmail(emailStr)
	if err != nil {
		return err
	}
	u.Email = email
	return nil
}

// UserRepository defines the persistence contract for Users
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*User, error)
	FindAll(ctx context.Context, limit, offset int) ([]User, int64, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}
