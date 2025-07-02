package models

import "time"

type User struct {
	ID          uint      `gorm:"primaryKey"`
	FirebaseUID string    `gorm:"uniqueIndex"`
	Email       string
	Name        string
	Role        string    `gorm:"default:user"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}