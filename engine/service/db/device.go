package db

import (
	"github.com/google/uuid"
)

type Device struct {
	Base
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
}

func (Device) TableName() string {
	return "device"
}
