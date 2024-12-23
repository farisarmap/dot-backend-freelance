package entity

import "time"

type Order struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderName string    `gorm:"size:100" json:"order_name"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User User `gorm:"foreignKey:UserID" json:"user"`
}
