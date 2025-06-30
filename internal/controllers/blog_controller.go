package controllers

import (
    "auth2_google/internal/models"
    "auth2_google/internal/services"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "io"
    "strings"
    "log"
)

type BlogController struct {
    blogService services.BlogServiceInterface
}

func NewBlogController(blogService services.BlogServiceInterface) *BlogController {
    return &BlogController{
        blogService: blogService,
    }
}

func (ctrl *BlogController) GetAllPosts(c *gin.Context) {
    posts, err := ctrl.blogService.GetAllPosts()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   "Failed to get posts",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "posts":   posts,
    })
}

func (ctrl *BlogController) GetPublishedPosts(c *gin.Context) {
    posts, err := ctrl.blogService.GetPublishedPosts()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   "Failed to get published posts",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "posts":   posts,
    })
}

func (ctrl *BlogController) GetPost(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid post ID",
        })
        return
    }

    post, err := ctrl.blogService.GetPostByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "error":   "Post not found",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "post":    post,
    })
}

func (ctrl *BlogController) CreatePost(c *gin.Context) {
    log.Println("=== INCOMING REQUEST DEBUG ===")
    log.Printf("Method: %s", c.Request.Method)
    log.Printf("URL: %s", c.Request.URL.String())
    log.Printf("Content-Type: %s", c.GetHeader("Content-Type"))
    
    // Read the raw body
    bodyBytes, _ := io.ReadAll(c.Request.Body)
    log.Printf("Raw Body: %s", string(bodyBytes))
    
    // Reset body for binding
    c.Request.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
    
    var req models.CreateBlogPostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        log.Printf("❌ Binding Error: %v", err)
        c.JSON(400, gin.H{
            "error": err.Error(),
            "received_body": string(bodyBytes),
        })
        return
    }
    
    log.Printf("✅ Parsed Successfully: %+v", req)
    post, err := ctrl.blogService.CreatePost(req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "success": true,
        "message": "Post created successfully",
        "post":    post,
    })
}

func (ctrl *BlogController) UpdatePost(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid post ID",
        })
        return
    }

    var req models.UpdateBlogPostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid input: " + err.Error(),
        })
        return
    }

    post, err := ctrl.blogService.UpdatePost(uint(id), req)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Post updated successfully",
        "post":    post,
    })
}

func (ctrl *BlogController) DeletePost(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid post ID",
        })
        return
    }

    err = ctrl.blogService.DeletePost(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Post deleted successfully",
    })
}