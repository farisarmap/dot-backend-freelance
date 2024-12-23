package adapter

import (
	"context"
	"database/sql"

	"github.com/farisarmap/dot-backend-freelance/internal/entity"
	"gorm.io/gorm"
)

type OrderRepository interface {
	GetAll(ctx context.Context) ([]entity.Order, error)
	GetByID(ctx context.Context, id uint) (entity.Order, error)
	Create(ctx context.Context, order *entity.Order) error
	Update(ctx context.Context, order *entity.Order) error
	Delete(ctx context.Context, order *entity.Order) error
	Transaction(ctx context.Context, fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) GetAll(ctx context.Context) ([]entity.Order, error) {
	var orders []entity.Order
	if err := r.db.WithContext(ctx).Preload("User").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) GetByID(ctx context.Context, id uint) (entity.Order, error) {
	var order entity.Order
	if err := r.db.WithContext(ctx).Preload("User").First(&order, id).Error; err != nil {
		return order, err
	}
	return order, nil
}

func (r *orderRepository) Create(ctx context.Context, order *entity.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderRepository) Update(ctx context.Context, order *entity.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

func (r *orderRepository) Delete(ctx context.Context, order *entity.Order) error {
	return r.db.WithContext(ctx).Delete(order).Error
}

func (r *orderRepository) Transaction(ctx context.Context, fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return r.db.WithContext(ctx).Transaction(fc, opts...)
}
