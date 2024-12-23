package service

import (
	"context"
	"encoding/json"

	"github.com/farisarmap/dot-backend-freelance/internal/adapter"
	"github.com/farisarmap/dot-backend-freelance/internal/entity"
	"gorm.io/gorm"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	GetUserByID(ctx context.Context, id uint) (entity.User, error)
	CreateUser(ctx context.Context, name, email string) (entity.User, error)
	UpdateUser(ctx context.Context, id uint, name, email string) (entity.User, error)
	PartialUpdateUser(ctx context.Context, id uint, name, email *string) (entity.User, error)
	DeleteUser(ctx context.Context, id uint) error
}

type userService struct {
	userRepo     adapter.UserRepository
	orderRepo    adapter.OrderRepository
	cacheManager adapter.CacheManager
}

func NewUserService(
	userRepo adapter.UserRepository,
	orderRepo adapter.OrderRepository,
	cacheManager adapter.CacheManager,
) UserService {
	return &userService{
		userRepo:     userRepo,
		orderRepo:    orderRepo,
		cacheManager: cacheManager,
	}
}

func (s *userService) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	cachedData, err := s.cacheManager.Get("users")
	if err == nil && cachedData != "" {
		var users []entity.User
		if err := json.Unmarshal([]byte(cachedData), &users); err == nil {
			return users, err
		}

		return users, nil
	}

	resp, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return []entity.User{}, err
	}
	for i := range resp {
		for j := range resp[i].Orders {
			resp[i].Orders[j].User = entity.User{}
		}
	}

	s.cacheManager.Set("users", resp)
	return resp, nil
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (entity.User, error) {
	cachedData, err := s.cacheManager.Get("users_detail")
	if err == nil && cachedData != "" {
		var users entity.User
		if err := json.Unmarshal([]byte(cachedData), &users); err == nil {
			return users, err
		}

		return users, nil
	}

	resp, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return entity.User{}, err
	}

	if err := s.cacheManager.Set("users_detail", resp); err != nil {
		return entity.User{}, err
	}

	return resp, nil
}

func (s *userService) CreateUser(ctx context.Context, name, email string) (entity.User, error) {
	if err := s.cacheManager.Delete("users"); err != nil {
		return entity.User{}, err
	}

	user := entity.User{
		Name:  name,
		Email: email,
	}

	err := s.userRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, id uint, name, email string) (entity.User, error) {
	var user entity.User

	err := s.userRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if err := tx.First(&user, id).Error; err != nil {
			return err
		}

		user.Name = name
		user.Email = email

		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		return nil
	})

	if err := s.cacheManager.Delete("users"); err != nil {
		return entity.User{}, err
	}

	if err := s.cacheManager.Delete("users_detail"); err != nil {
		return entity.User{}, err
	}

	return user, err
}

func (s *userService) PartialUpdateUser(ctx context.Context, id uint, name, email *string) (entity.User, error) {
	var user entity.User

	err := s.userRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if err := tx.First(&user, id).Error; err != nil {
			return err
		}

		if name != nil {
			user.Name = *name
		}
		if email != nil {
			user.Email = *email
		}

		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		return nil
	})

	if err := s.cacheManager.Delete("users"); err != nil {
		return entity.User{}, err
	}

	if err := s.cacheManager.Delete("users_detail"); err != nil {
		return entity.User{}, err
	}

	return user, err
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.cacheManager.Delete("users"); err != nil {
		return err
	}

	if err := s.cacheManager.Delete("users_detail"); err != nil {
		return err
	}

	return s.userRepo.Delete(ctx, &user)
}
