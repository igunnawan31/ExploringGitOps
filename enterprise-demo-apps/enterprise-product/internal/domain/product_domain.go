package domain

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        string `gorm:"primaryKey;type:uuid"`
	Name      string `gorm:"type:varchar(100);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Product struct {
	ID          string   `gorm:"primaryKey;type:uuid"`
	CategoryID  string   `gorm:"type:uuid;index"`
	Category    Category `gorm:"foreignKey:CategoryID"`
	Name        string   `gorm:"type:varchar(200);not null"`
	SKU         string   `gorm:"type:varchar(50);uniqueIndex;not null"`
	Description string   `gorm:"type:text"`
	Price       float64  `gorm:"type:decimal(10,2);not null"`
	Stock       int      `gorm:"type:int;not null;default:0"`
	ImageURL    string   `gorm:"type:varchar(255)"`
	Status      string   `gorm:"type:varchar(20);default:'active'"` // active, draft
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
