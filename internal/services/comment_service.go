package services

import (
    "auth2_google/internal/models"
    "auth2_google/internal/repositories"
    "errors"
    "time"
)

type CommentServiceInterface interface {
    CreateComment(req models.CreateCommentRequest) (*models.CommentResponse, error)
    GetCommentsByBlogPostID(blogPostID uint) ([]models.CommentResponse, error)
    GetCommentByID(id uint) (*models.CommentResponse, error)
    UpdateComment(id uint, req models.UpdateCommentRequest) (*models.CommentResponse, error)
    DeleteComment(id uint) error
    CreateReply(parentID uint, req models.CreateCommentRequest) (*models.CommentResponse, error)
    GetReplies(parentID uint) ([]models.CommentResponse, error)
}

type CommentService struct {
    commentRepo repositories.CommentRepositoryInterface
}

func NewCommentService(commentRepo repositories.CommentRepositoryInterface) CommentServiceInterface {
    return &CommentService{
        commentRepo: commentRepo,
    }
}

// Helper function to format date
func formatCommentDate(t time.Time) string {
    return t.Format("January 2, 2006 at 3:04 PM") // "June 23, 2025 at 4:30 PM"
}

// Convert to response format with nested replies
func (s *CommentService) toResponse(comment models.Comment) models.CommentResponse {
    response := models.CommentResponse{
        ID:         comment.ID,
        BlogPostID: comment.BlogPostID,
        Name:       comment.Name,
        Email:      comment.Email,
        Text:       comment.Text,
        ParentID:   comment.ParentID,
        CreatedAt:  formatCommentDate(comment.CreatedAt),
        Replies:    []models.CommentResponse{},
    }

    // Convert replies to response format
    for _, reply := range comment.Replies {
        replyResponse := s.toResponse(reply) // Recursive for nested replies
        response.Replies = append(response.Replies, replyResponse)
    }

    return response
}

func (s *CommentService) CreateComment(req models.CreateCommentRequest) (*models.CommentResponse, error) {
    if req.Name == "" || req.Text == "" {
        return nil, errors.New("name and text are required")
    }

    comment := &models.Comment{
        BlogPostID: req.BlogPostID,
        Name:       req.Name,
        Email:      req.Email,
        Text:       req.Text,
        ParentID:   req.ParentID,
    }

    err := s.commentRepo.Create(comment)
    if err != nil {
        return nil, err
    }

    response := s.toResponse(*comment)
    return &response, nil
}

func (s *CommentService) GetCommentsByBlogPostID(blogPostID uint) ([]models.CommentResponse, error) {
    comments, err := s.commentRepo.GetByBlogPostID(blogPostID)
    if err != nil {
        return nil, err
    }

    var responses []models.CommentResponse
    for _, comment := range comments {
        responses = append(responses, s.toResponse(comment))
    }

    return responses, nil
}

func (s *CommentService) GetCommentByID(id uint) (*models.CommentResponse, error) {
    comment, err := s.commentRepo.GetByID(id)
    if err != nil {
        return nil, errors.New("comment not found")
    }

    response := s.toResponse(*comment)
    return &response, nil
}

func (s *CommentService) UpdateComment(id uint, req models.UpdateCommentRequest) (*models.CommentResponse, error) {
    comment, err := s.commentRepo.GetByID(id)
    if err != nil {
        return nil, errors.New("comment not found")
    }

    if req.Text != nil {
        comment.Text = *req.Text
    }

    err = s.commentRepo.Update(comment)
    if err != nil {
        return nil, err
    }

    response := s.toResponse(*comment)
    return &response, nil
}

func (s *CommentService) DeleteComment(id uint) error {
    _, err := s.commentRepo.GetByID(id)
    if err != nil {
        return errors.New("comment not found")
    }

    return s.commentRepo.Delete(id)
}

func (s *CommentService) CreateReply(parentID uint, req models.CreateCommentRequest) (*models.CommentResponse, error) {
    if req.Name == "" || req.Text == "" {
        return nil, errors.New("name and text are required")
    }

    // Verify parent comment exists
    parentComment, err := s.commentRepo.GetByID(parentID)
    if err != nil {
        return nil, errors.New("parent comment not found")
    }

    reply := &models.Comment{
        BlogPostID: parentComment.BlogPostID, // Same blog as parent
        Name:       req.Name,
        Email:      req.Email,
        Text:       req.Text,
        ParentID:   &parentID, // Set parent ID
    }

    err = s.commentRepo.Create(reply)
    if err != nil {
        return nil, err
    }

    response := s.toResponse(*reply)
    return &response, nil
}

func (s *CommentService) GetReplies(parentID uint) ([]models.CommentResponse, error) {
    replies, err := s.commentRepo.GetReplies(parentID)
    if err != nil {
        return nil, err
    }

    var responses []models.CommentResponse
    for _, reply := range replies {
        responses = append(responses, s.toResponse(reply))
    }

    return responses, nil
}