package controllers

import (
    "auth2_google/internal/services"
    "auth2_google/internal/models"
    "github.com/gin-gonic/gin"
    "net/http"
    "strconv"
)

type CommentController struct {
    commentService services.CommentServiceInterface
}

func NewCommentController(commentService services.CommentServiceInterface) *CommentController {
    return &CommentController{
        commentService: commentService,
    }
}

// POST /api/blogs/:id/comments - Create comment for a blog
func (ctrl *CommentController) CreateComment(c *gin.Context) {
    blogIDStr := c.Param("id")
    blogID, err := strconv.ParseUint(blogIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid blog ID",
        })
        return
    }

    var req models.CreateCommentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid input: " + err.Error(),
        })
        return
    }

    // Set the blog post ID from URL parameter
    req.BlogPostID = uint(blogID)

    comment, err := ctrl.commentService.CreateComment(req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "success": true,
        "message": "Comment created successfully",
        "comment": comment,
    })
}

// GET /api/blogs/:id/comments - Get all comments for a blog
func (ctrl *CommentController) GetCommentsByBlog(c *gin.Context) {
    blogIDStr := c.Param("id")
    blogID, err := strconv.ParseUint(blogIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid blog ID",
        })
        return
    }

    comments, err := ctrl.commentService.GetCommentsByBlogPostID(uint(blogID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   "Failed to get comments",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success":  true,
        "comments": comments,
    })
}

// GET /api/comments/:id - Get single comment with replies
func (ctrl *CommentController) GetComment(c *gin.Context) {
    commentIDStr := c.Param("id")
    commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid comment ID",
        })
        return
    }

    comment, err := ctrl.commentService.GetCommentByID(uint(commentID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "comment": comment,
    })
}

// PUT /api/comments/:id - Update comment
func (ctrl *CommentController) UpdateComment(c *gin.Context) {
    commentIDStr := c.Param("id")
    commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid comment ID",
        })
        return
    }

    var req models.UpdateCommentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid input: " + err.Error(),
        })
        return
    }

    comment, err := ctrl.commentService.UpdateComment(uint(commentID), req)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Comment updated successfully",
        "comment": comment,
    })
}

// DELETE /api/comments/:id - Delete comment
func (ctrl *CommentController) DeleteComment(c *gin.Context) {
    commentIDStr := c.Param("id")
    commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid comment ID",
        })
        return
    }

    err = ctrl.commentService.DeleteComment(uint(commentID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Comment deleted successfully",
    })
}

// POST /api/comments/:id/reply - Create reply to a comment
func (ctrl *CommentController) CreateReply(c *gin.Context) {
    parentIDStr := c.Param("id")
    parentID, err := strconv.ParseUint(parentIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid parent comment ID",
        })
        return
    }

    var req models.CreateCommentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid input: " + err.Error(),
        })
        return
    }

    reply, err := ctrl.commentService.CreateReply(uint(parentID), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "success": true,
        "message": "Reply created successfully",
        "reply":   reply,
    })
}

// GET /api/comments/:id/replies - Get all replies for a comment
func (ctrl *CommentController) GetReplies(c *gin.Context) {
    parentIDStr := c.Param("id")
    parentID, err := strconv.ParseUint(parentIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid parent comment ID",
        })
        return
    }

    replies, err := ctrl.commentService.GetReplies(uint(parentID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   "Failed to get replies",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "replies": replies,
    })
}