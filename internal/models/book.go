package models

import "time"

type Book struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Title     string `json:"title" validate:"required"`
	Author    string `json:"author" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
