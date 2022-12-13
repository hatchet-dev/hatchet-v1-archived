package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"gorm.io/gorm"
)

type Base struct {
	ID        string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`
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

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
