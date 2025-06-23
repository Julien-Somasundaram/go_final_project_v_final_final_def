package models

import (
	"time"
)

// Link représente un lien raccourci dans la base de données.
type Link struct {
	ID        uint   `gorm:"primaryKey"`
	ShortCode string `gorm:"column:short_code;type:varchar(10);uniqueIndex;not null"`
	LongURL   string `gorm:"type:text;not null"`
	CreatedAt time.Time
}
