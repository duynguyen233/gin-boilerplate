package repositories

import (
	"context"
	"template/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user model.User) (model.User, error) {
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", user.ID).Updates(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.User{}).Error
}
