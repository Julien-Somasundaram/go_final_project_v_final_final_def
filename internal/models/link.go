package models

import "time"

type Link struct {
	ID        uint   `gorm:"primaryKey"`
	ShortCode string `gorm:"column:shortcode;type:varchar(10);uniqueIndex;not null"`
	LongURL   string `gorm:"type:text;not null"`
	CreatedAt time.Time
}
