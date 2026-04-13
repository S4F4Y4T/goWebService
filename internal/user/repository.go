package user

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// userSchema is the Persistence Model (Database Schema)
type userSchema struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:100"`
	Email     string    `gorm:"uniqueIndex;size:255"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (userSchema) TableName() string {
	return "users"
}

// fromDomain converts a Domain Entity to a Persistence Model
func fromDomain(u *User) userSchema {
	return userSchema{
		ID:        u.ID,
		Name:      u.Name,
		Email:     string(u.Email),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// toDomain converts a Persistence Model back to a Domain Entity
func (s userSchema) toDomain() *User {
	return &User{
		ID:        s.ID,
		Name:      s.Name,
		Email:     Email(s.Email),
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, u *User) error {
	schema := fromDomain(u)
	if err := r.db.WithContext(ctx).Create(&schema).Error; err != nil {
		return err
	}
	u.ID = schema.ID
	u.CreatedAt = schema.CreatedAt
	u.UpdatedAt = schema.UpdatedAt
	return nil
}

func (r *userRepository) Update(ctx context.Context, u *User) error {
	schema := fromDomain(u)
	return r.db.WithContext(ctx).Save(&schema).Error
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var s userSchema
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return s.toDomain(), nil
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&userSchema{}, id).Error
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*User, error) {
	var s userSchema
	if err := r.db.WithContext(ctx).First(&s, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return s.toDomain(), nil
}

func (r *userRepository) FindAll(ctx context.Context, limit, offset int) ([]User, int64, error) {
	var schemas []userSchema
	var total int64

	query := r.db.WithContext(ctx).Model(&userSchema{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&schemas).Error; err != nil {
		return nil, 0, err
	}

	users := make([]User, len(schemas))
	for i, s := range schemas {
		users[i] = *s.toDomain()
	}

	return users, total, nil
}
