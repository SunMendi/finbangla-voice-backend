package repositories

import (
	 "auth2_google/internal/models"
	 "gorm.io/gorm"
)


type BlogRepositoryInterface interface {
	 Create(post *models.BlogPost) error 
	 GetAll() ([]models.BlogPost, error)
	 GetByID(id uint) (*models.BlogPost, error)
	 Update(post *models.BlogPost) error 
	 Delete(id uint) error
	 GetPublished() ([]models.BlogPost, error) 
}

type blogRepository struct {
	db *gorm.DB 
}

func NewBlogRepository(db *gorm.DB) BlogRepositoryInterface {
	return &blogRepository {
		db: db,
	}
}

func(r *blogRepository) Create(post *models.BlogPost) error {
	return r.db.Create(post).Error 
}
func (r *blogRepository) GetAll() ([]models.BlogPost, error) {
    var posts []models.BlogPost
    err := r.db.Order("created_at DESC").Find(&posts).Error
    return posts, err
}
func (r *blogRepository) GetByID(id uint) (*models.BlogPost, error) {
    var post models.BlogPost
    err := r.db.First(&post, id).Error
    if err != nil {
        return nil, err
    }
    return &post, nil
}

func (r *blogRepository) Update(post *models.BlogPost) error {
    return r.db.Save(post).Error
}

func (r *blogRepository) Delete(id uint) error {
    return r.db.Delete(&models.BlogPost{}, id).Error
}

func (r *blogRepository) GetPublished() ([]models.BlogPost, error) {
    var posts []models.BlogPost
    err := r.db.Where("published = ?", true).Order("created_at DESC").Find(&posts).Error
    return posts, err
}