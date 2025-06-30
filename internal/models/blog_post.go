package models

import (
	"time"
	"gorm.io/gorm"
)
type BlogPost struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Title     string         `json:"title" gorm:"not null"`
    Excerpt   string         `json:"excerpt" gorm:"type:text;not null"` // 🔥 NEW
    Author    string         `json:"author" gorm:"not null"`            // 🔥 NEW
    Image     string         `json:"image"`                             // 🔥 NEW
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

    Comments  []Comment       `json:"comments,omitempty" gorm:"foreignKey:BlogPostID"`
}

// Request DTOs
type CreateBlogPostRequest struct {
    Title   string `json:"title" binding:"required"`
    Excerpt string `json:"excerpt" binding:"required"` // 🔥 NEW
    Author  string `json:"author" binding:"required"`  // 🔥 NEW
    Image   string `json:"image"`                       // 🔥 NEW
}

type UpdateBlogPostRequest struct {
    Title   *string `json:"title"`
    Excerpt *string `json:"excerpt"` // 🔥 NEW
    Author  *string `json:"author"`  // 🔥 NEW
    Image   *string `json:"image"`   // 🔥 NEW
}

//Response Data Transfer Model 

type BlogPostResponse struct {
    ID      string `json:"id"`      // 🔥 String for frontend
    Title   string `json:"title"`
    Excerpt string `json:"excerpt"` // 🔥 NEW
    Author  string `json:"author"`  // 🔥 NEW
    Date    string `json:"date"`    // 🔥 NEW - Formatted date
    Image   string `json:"image"`   // 🔥 NEW
}