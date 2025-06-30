package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	 ID    uint  `json:"id" gorm:"primaryKey"`
	 GoogleID string `json:"google_id" gorm:"uniqueIndex;not null"`
	 Email string `json:"email" gorm:"uniqueIndex;not null"`
	 Name string `json:"name" gorm:"not null"`
	 Picture   string    `json:"picture"`
	 CreatedAt time.Time  `json:"created_at"`
	 UpdatedAt time.Time `json:"updated_at"`
	 DeletedAt gorm.DeletedAt `json:"_" gorm:"index"`
}

// This is what Google will send us
type GoogleUserInfo struct {
    ID      string `json:"id"`      // Google's user ID
    Email   string `json:"email"`   // user@gmail.com
    Name    string `json:"name"`    // "John Doe"
    Picture string `json:"picture"` // Profile photo URL
}