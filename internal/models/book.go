package models

import "time"

type Book struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title" validate:"required,min=2,max=200"`
	Author    string    `json:"author" validate:"required,min=2,max=100"`
	ISBN      string    `json:"isbn" validate:"omitempty,len=13"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
