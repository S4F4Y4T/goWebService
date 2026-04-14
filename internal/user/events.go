package user

import (
	"time"

	"github.com/S4F4Y4T/goWebService/internal/shared/domain"
)

const UserCreatedTopic = "user.created"

// UserCreated is emitted when a new user is created
type UserCreated struct {
	domain.BaseEvent
	UserID uint
	Email  string
}

func NewUserCreated(userID uint, email string) UserCreated {
	return UserCreated{
		BaseEvent: domain.BaseEvent{Timestamp: time.Now()},
		UserID:    userID,
		Email:     email,
	}
}

func (e UserCreated) Topic() string {
	return UserCreatedTopic
}
