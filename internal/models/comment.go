package models 

import (
    "time"
    "gorm.io/gorm"
)

type Comment struct {
    ID         uint           `json:"id" gorm:"primaryKey"`
    BlogPostID uint           `json:"blog_post_id" gorm:"not null;index"`
    Name       string         `json:"name" gorm:"not null"`
    Email      string         `json:"email"`                    // Fixed: lowercase 'email'
    Text       string         `json:"text" gorm:"type:text;not null"` // 🔥 MISSING - Add this field
    ParentID   *uint          `json:"parent_id" gorm:"index"`
    CreatedAt  time.Time      `json:"created_at"`
    UpdatedAt  time.Time      `json:"updated_at"`
    DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

    // Relationships
    BlogPost BlogPost  `json:"-" gorm:"foreignKey:BlogPostID"`
    Parent   *Comment  `json:"-" gorm:"foreignKey:ParentID"`
    Replies  []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
}

type CreateCommentRequest struct {
    BlogPostID uint   `json:"blog_post_id"`
    Name       string `json:"name" binding:"required"`
    Email      string `json:"email"`
    Text       string `json:"text" binding:"required"`
    ParentID   *uint  `json:"parent_id"` // null for main comments, ID for replies
}

type UpdateCommentRequest struct {
    Text *string `json:"text"`
}

type CommentResponse struct {
    ID         uint              `json:"id"`
    BlogPostID uint              `json:"blog_post_id"`
    Name       string            `json:"name"`
    Email      string            `json:"email,omitempty"`
    Text       string            `json:"text"`
    ParentID   *uint             `json:"parent_id"`
    CreatedAt  string            `json:"created_at"` // Formatted date
    Replies    []CommentResponse `json:"replies,omitempty"`
}