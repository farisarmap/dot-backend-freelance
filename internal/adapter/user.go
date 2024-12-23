package adapter

import (
	"context"
	"database/sql"

	"github.com/farisarmap/dot-backend-freelance/internal/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]entity.User, error)
	GetByID(ctx context.Context, id uint) (entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, user *entity.User) error
	Transaction(ctx context.Context, fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetAll(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	if err := r.db.WithContext(ctx).Preload("Orders").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Preload("Orders").First(&user, id).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Delete(user).Error
}

func (r *userRepository) Transaction(ctx context.Context, fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return r.db.WithContext(ctx).Transaction(fc, opts...)
}
