package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"gorm.io/gorm"
)

type Base struct {
	ID        string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time      `gorm:"index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}

	return nil
}

func (b *Base) ToAPITypeMetadata() *types.APIResourceMeta {
	return &types.APIResourceMeta{
		CreatedAt: &b.CreatedAt,
		UpdatedAt: &b.UpdatedAt,
		ID:        b.ID,
	}
}
