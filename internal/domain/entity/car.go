package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Car struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name          string         `json:"name" gorm:"type:varchar(100);not null;index:idx_cars_name"`
	EngineVersion string         `json:"engine_version" gorm:"type:varchar(10);not null;index:idx_cars_engine_version"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime;index:idx_cars_created_at"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index:idx_cars_deleted_at"`
}

func (Car) TableName() string {
	return "cars"
}

// BeforeCreate hook to generate UUID before creating
func (c *Car) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
