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

// HasEncryptedFields is used for models which have an encrypted field.
// After Encrypt() and Decrypt() methods are called, these methods should set FieldsAreEncrypted
// correspondingly.
type HasEncryptedFields struct {
	FieldsAreEncrypted bool `gorm:"-"`
}

func (h *HasEncryptedFields) AfterFind(tx *gorm.DB) (err error) {
	h.FieldsAreEncrypted = true

	return
}
