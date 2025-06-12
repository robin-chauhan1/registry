package models

import "time"

type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" binding:"required"`
    Email     string    `json:"email" binding:"required,email" gorm:"unique"`
    Password  string    `json:"-" binding:"required"` // Use `-` to exclude from JSON responses
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}