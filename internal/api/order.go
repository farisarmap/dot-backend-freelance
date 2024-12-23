package api

type CreateOrder struct {
	OrderName string `json:"order_name" validate:"required,min=3,max=100"`
	UserID    uint   `json:"user_id" validate:"required"`
}

type PartiallyUpdateOrder struct {
	OrderName *string `json:"order_name" validate:"omitempty,min=3,max=100"`
	UserID    *uint   `json:"user_id"`
}

type CreateUserAndOrderRequest struct {
	User  CreateUserRequest  `json:"user" validate:"required"`
	Order CreateOrderRequest `json:"order" validate:"required"`
}

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type CreateOrderRequest struct {
	OrderName string `json:"order_name" validate:"required"`
}
