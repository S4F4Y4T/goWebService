package auth

import "time"

// Credential is the domain model for authentication credentials.
type Credential struct {
	ID           uint
	UserID       uint
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

// CredentialRepository defines the persistence contract.
type CredentialRepository interface {
	Save(cred *Credential) error
	FindByEmail(email string) (*Credential, error)
}
