package repository

import (
	"context"
	"errors"

	"github.com/S4F4Y4T/goWebService/internal/model"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) model.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	user := model.User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, req *model.UpdateUserRequest) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.Name = req.Name
	user.Email = req.Email
	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil if not found
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Delete(ctx context.Context, req *model.DeleteUserRequest) error {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if err := r.db.WithContext(ctx).Delete(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, req *model.GetUserRequest) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindAll(ctx context.Context, req *model.GetUsersRequest) (*model.GetUsersResponse, error) {
	var users []model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{})
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if req.Limit > 0 {
		query = query.Limit(req.Limit)
	}
	if req.Offset > 0 {
		query = query.Offset(req.Offset)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return &model.GetUsersResponse{
		Users:  users,
		Total:  total,
		Limit:  req.Limit,
		Offset: req.Offset,
	}, nil
}
