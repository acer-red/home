package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type Base struct {
	ID        uuid.UUID `gorm:"primarykey;type:uuid;default:uuidv7()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate hook to generate UUIDv7 if not present
func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID, err = uuid.NewV7()
	}
	return
}

func Init(dsn string) error {
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// AutoMigrate
	return db.AutoMigrate(&User{}, &API{}, &File{}, &Feedback{}, &Product{}, &Device{})
}

func Disconnect() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
