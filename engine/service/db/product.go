package db

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID        uuid.UUID `gorm:"primarykey;type:uuid;default:uuidv7()" json:"id"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}

func (Product) TableName() string {
	return "product"
}
