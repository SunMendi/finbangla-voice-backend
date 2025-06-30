package services

import (
    "auth2_google/internal/models"
    "auth2_google/internal/repositories"
    "errors"
    "fmt"
    "time"
)

type BlogServiceInterface interface {
    CreatePost(req models.CreateBlogPostRequest) (*models.BlogPostResponse, error)
    GetAllPosts() ([]models.BlogPostResponse, error)
    GetPostByID(id uint) (*models.BlogPostResponse, error)
    UpdatePost(id uint, req models.UpdateBlogPostRequest) (*models.BlogPostResponse, error)
    DeletePost(id uint) error
    GetPublishedPosts() ([]models.BlogPostResponse, error) // Keep existing
}

type BlogService struct {
    blogRepo repositories.BlogRepositoryInterface
}

func NewBlogService(blogRepo repositories.BlogRepositoryInterface) BlogServiceInterface {
    return &BlogService{
        blogRepo: blogRepo,
    }
}

// 🔥 Helper function to format date
func formatDate(t time.Time) string {
    return t.Format("January 2, 2006") // "June 20, 2025"
}

// 🔥 Convert to response format
func (s *BlogService) toResponse(post models.BlogPost) models.BlogPostResponse {
    return models.BlogPostResponse{
        ID:      fmt.Sprintf("%d", post.ID), // Convert to string
        Title:   post.Title,
        Excerpt: post.Excerpt, // 🔥 NEW
        Author:  post.Author,  // 🔥 NEW
        Date:    formatDate(post.CreatedAt),
        Image:   post.Image,   // 🔥 NEW
    }
}

func (s *BlogService) CreatePost(req models.CreateBlogPostRequest) (*models.BlogPostResponse, error) {
    if req.Title == "" || req.Excerpt == "" || req.Author == "" {
        return nil, errors.New("title, excerpt and author are required")
    }

    post := &models.BlogPost{
        Title:   req.Title,
        Excerpt: req.Excerpt, // 🔥 NEW
        Author:  req.Author,  // 🔥 NEW
        Image:   req.Image,   // 🔥 NEW
    }

    err := s.blogRepo.Create(post)
    if err != nil {
        return nil, err
    }

    response := s.toResponse(*post)
    return &response, nil
}

func (s *BlogService) GetAllPosts() ([]models.BlogPostResponse, error) {
    posts, err := s.blogRepo.GetAll()
    if err != nil {
        return nil, err
    }

    var responses []models.BlogPostResponse
    for _, post := range posts {
        responses = append(responses, s.toResponse(post))
    }

    return responses, nil
}

func (s *BlogService) GetPostByID(id uint) (*models.BlogPostResponse, error) {
    post, err := s.blogRepo.GetByID(id)
    if err != nil {
        return nil, errors.New("post not found")
    }

    response := s.toResponse(*post)
    return &response, nil
}

func (s *BlogService) UpdatePost(id uint, req models.UpdateBlogPostRequest) (*models.BlogPostResponse, error) {
    post, err := s.blogRepo.GetByID(id)
    if err != nil {
        return nil, errors.New("post not found")
    }

    if req.Title != nil {
        post.Title = *req.Title
    }
    if req.Excerpt != nil {
        post.Excerpt = *req.Excerpt // 🔥 NEW
    }
    if req.Author != nil {
        post.Author = *req.Author // 🔥 NEW
    }
    if req.Image != nil {
        post.Image = *req.Image // 🔥 NEW
    }

    err = s.blogRepo.Update(post)
    if err != nil {
        return nil, err
    }

    response := s.toResponse(*post)
    return &response, nil
}

func (s *BlogService) DeletePost(id uint) error {
    _, err := s.blogRepo.GetByID(id)
    if err != nil {
        return errors.New("post not found")
    }

    return s.blogRepo.Delete(id)
}

func (s *BlogService) GetPublishedPosts() ([]models.BlogPostResponse, error) {
    posts, err := s.blogRepo.GetPublished()
    if err != nil {
        return nil, err
    }

    var responses []models.BlogPostResponse
    for _, post := range posts {
        responses = append(responses, s.toResponse(post))
    }

    return responses, nil
}