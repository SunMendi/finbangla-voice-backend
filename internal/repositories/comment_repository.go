package repositories

import (
    "auth2_google/internal/models"
    "gorm.io/gorm"
)

type CommentRepositoryInterface interface {
    Create(comment *models.Comment) error
    GetByBlogPostID(blogPostID uint) ([]models.Comment, error)
    GetByID(id uint) (*models.Comment, error)
    Update(comment *models.Comment) error
    Delete(id uint) error
    GetReplies(parentID uint) ([]models.Comment, error)
}

type CommentRepository struct {
    db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepositoryInterface {
    return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment *models.Comment) error {
    return r.db.Create(comment).Error
}

func (r *CommentRepository) GetByBlogPostID(blogPostID uint) ([]models.Comment, error) {
    var comments []models.Comment
    
    // Get only main comments (ParentID is null) with their replies
    err := r.db.Where("blog_post_id = ? AND parent_id IS NULL", blogPostID).
        Preload("Replies").
        Order("created_at DESC").
        Find(&comments).Error
    
    return comments, err
}

func (r *CommentRepository) GetByID(id uint) (*models.Comment, error) {
    var comment models.Comment
    err := r.db.Preload("Replies").First(&comment, id).Error
    return &comment, err
}

func (r *CommentRepository) Update(comment *models.Comment) error {
    return r.db.Save(comment).Error
}

func (r *CommentRepository) Delete(id uint) error {
    return r.db.Delete(&models.Comment{}, id).Error
}

func (r *CommentRepository) GetReplies(parentID uint) ([]models.Comment, error) {
    var replies []models.Comment
    err := r.db.Where("parent_id = ?", parentID).
        Order("created_at ASC").
        Find(&replies).Error
    return replies, err
}