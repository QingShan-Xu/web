package bm

import (
	"time"

	"gorm.io/gorm"
)

type ID struct {
	ID        int            `gorm:"primarykey,type:integer(11)" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
