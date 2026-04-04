package repository

import (
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

func (r *userRepository) Create(req *model.CreateUserRequest) (*model.User, error) {
	user := model.User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Update(req *model.UpdateUserRequest) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.Name = req.Name
	if err := r.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Delete(req *model.DeleteUserRequest) error {
	var user model.User
	if err := r.db.First(&user, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if err := r.db.Delete(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) FindByID(req *model.GetUserRequest) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindAll(req *model.GetUsersRequest) (*model.GetUsersResponse, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{})
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
