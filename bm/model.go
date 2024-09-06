package bm

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        int            `gorm:"primarykey,type:integer(11)" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"delete_at"`
}
