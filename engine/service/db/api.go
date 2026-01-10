package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type API struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	APIKey     string         `json:"api_key" gorm:"unique"`
	ExpiresAt  time.Time      `json:"expires_at"`
	LastUsedAt time.Time      `json:"last_used_at"`
	UsedTimes  int32          `json:"used_times"`
	ProductID  uuid.UUID      `json:"product_id" gorm:"type:uuid;not null"`
}

// 设置表名
func (API) TableName() string {
	return "api"
}
