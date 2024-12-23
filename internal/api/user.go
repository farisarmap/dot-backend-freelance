package api

type CreateUser struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type PartiallyUpdateUser struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
}
