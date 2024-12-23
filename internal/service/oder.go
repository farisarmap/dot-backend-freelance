package service

import (
	"context"
	"errors"

	"github.com/farisarmap/dot-backend-freelance/internal/adapter"
	"github.com/farisarmap/dot-backend-freelance/internal/api"
	"github.com/farisarmap/dot-backend-freelance/internal/entity"
	"gorm.io/gorm"
)

type OrderService interface {
	GetAllOrders(ctx context.Context) ([]entity.Order, error)
	GetOrderByID(ctx context.Context, id uint) (entity.Order, error)
	CreateOrder(ctx context.Context, orderName string, userID uint) (entity.Order, error)
	UpdateOrder(ctx context.Context, id uint, orderName string, userID uint) (entity.Order, error)
	PartialUpdateOrder(ctx context.Context, id uint, orderName *string, userID *uint) (entity.Order, error)
	DeleteOrder(ctx context.Context, id uint) error

	CreateUserAndOrder(ctx context.Context, req api.CreateUserAndOrderRequest) error
}

type orderService struct {
	orderRepo adapter.OrderRepository
	userRepo  adapter.UserRepository
}

func NewOrderService(
	orderRepo adapter.OrderRepository,
	userRepo adapter.UserRepository,
) OrderService {
	return &orderService{
		orderRepo: orderRepo,
		userRepo:  userRepo,
	}
}

func (s *orderService) GetAllOrders(ctx context.Context) ([]entity.Order, error) {
	resp, err := s.orderRepo.GetAll(ctx)
	if err != nil {
		return []entity.Order{}, err
	}
	return resp, nil
}

func (s *orderService) GetOrderByID(ctx context.Context, id uint) (entity.Order, error) {
	return s.orderRepo.GetByID(ctx, id)
}

func (s *orderService) CreateOrder(ctx context.Context, orderName string, userID uint) (entity.Order, error) {
	order := entity.Order{
		OrderName: orderName,
		UserID:    userID,
	}

	if err := s.orderRepo.Transaction(ctx, func(tx *gorm.DB) error {
		err := tx.Save(&order).Error
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return entity.Order{}, err
	}

	return order, nil
}

func (s *orderService) UpdateOrder(ctx context.Context, id uint, orderName string, userID uint) (entity.Order, error) {
	var order entity.Order
	if err := s.orderRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if err := tx.First(&order, id).Error; err != nil {
			return err
		}

		order.OrderName = orderName
		order.UserID = userID

		if err := tx.Save(&order).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return entity.Order{}, err
	}

	return order, nil
}

func (s *orderService) PartialUpdateOrder(ctx context.Context, id uint, orderName *string, userID *uint) (entity.Order, error) {
	var order entity.Order
	if err := s.orderRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if err := tx.First(&order, id).Error; err != nil {
			return err
		}
		if orderName != nil {
			order.OrderName = *orderName
		}
		if userID != nil {
			_, err := s.userRepo.GetByID(ctx, *userID)
			if err != nil {
				return errors.New("user tidak ditemukan")
			}
			order.UserID = *userID
		}

		if err := tx.Save(&order).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return entity.Order{}, err
	}

	return order, nil
}

func (s *orderService) DeleteOrder(ctx context.Context, id uint) error {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.orderRepo.Delete(ctx, &order)
}

func (s *orderService) CreateUserAndOrder(ctx context.Context, req api.CreateUserAndOrderRequest) error {
	return s.orderRepo.Transaction(ctx, func(tx *gorm.DB) error {
		user := entity.User{
			Name:  req.User.Name,
			Email: req.User.Email,
		}

		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		order := entity.Order{
			OrderName: req.Order.OrderName,
			UserID:    user.ID,
		}

		if err := tx.Save(&order).Error; err != nil {
			return err
		}
		return nil
	})
}
