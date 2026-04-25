package auth

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// credentialSchema is the GORM persistence model.
type credentialSchema struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"uniqueIndex;not null"`
	Email        string    `gorm:"uniqueIndex;size:255;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (credentialSchema) TableName() string { return "credentials" }

type credentialRepository struct {
	db *gorm.DB
}

func NewCredentialRepository(db *gorm.DB) CredentialRepository {
	return &credentialRepository{db: db}
}

func (r *credentialRepository) Save(cred *Credential) error {
	schema := credentialSchema{
		UserID:       cred.UserID,
		Email:        cred.Email,
		PasswordHash: cred.PasswordHash,
	}
	if err := r.db.Create(&schema).Error; err != nil {
		return err
	}
	cred.ID = schema.ID
	cred.CreatedAt = schema.CreatedAt
	return nil
}

func (r *credentialRepository) FindByEmail(email string) (*Credential, error) {
	var s credentialSchema
	if err := r.db.Where("email = ?", email).First(&s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &Credential{
		ID:           s.ID,
		UserID:       s.UserID,
		Email:        s.Email,
		PasswordHash: s.PasswordHash,
		CreatedAt:    s.CreatedAt,
	}, nil
}
