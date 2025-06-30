package main

import (
    "auth2_google/internal/config"
    "auth2_google/internal/controllers"
    "auth2_google/internal/middleware"
    "auth2_google/internal/models"
    "auth2_google/internal/repositories"
    "auth2_google/internal/services"
    "auth2_google/pkg/database"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables (only works locally)
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found - using system environment variables")
    }

    // 🔥 EXPLICITLY set Gin mode based on GIN_MODE environment variable
    if os.Getenv("GIN_MODE") == "release" {
        gin.SetMode(gin.ReleaseMode)
        log.Println("🚀 Starting application in RELEASE mode")
    } else {
        gin.SetMode(gin.DebugMode)
        log.Println("🔧 Starting application in DEBUG mode")
    }

    // 🔥 Check if all required environment variables exist
    validateEnvironmentVariables()

    // Connect to database
    database.ConnectDatabase()

    // Auto-migrate database tables
    database.DB.AutoMigrate(&models.User{}, &models.BlogPost{}, &models.Comment{})
    log.Println("✅ Database tables created/updated")

    // Initialize Google OAuth2 configuration
    config.InitGoogleOAuth()
    log.Println("✅ Google OAuth2 configured")

    // Dependency injection
    blogRepo := repositories.NewBlogRepository(database.DB)
    blogService := services.NewBlogService(blogRepo)
    blogController := controllers.NewBlogController(blogService)

    commentRepo := repositories.NewCommentRepository(database.DB)
    commentService := services.NewCommentService(commentRepo)
    commentController := controllers.NewCommentController(commentService)

    // Setup Gin router
    router := gin.New() // Use gin.New() for more control over middleware

    // Add middleware
    router.Use(middleware.BasicLogger())
    router.Use(gin.Recovery()) // Handle panics gracefully

    // CORS configuration
    router.Use(cors.New(cors.Config{
        AllowOrigins: []string{
        "https://finbanglavoice.fi",      // 🔥 Your production domain
        "https://www.finbanglavoice.fi",  // 🔥 With www subdomain
        "http://localhost:3000",          // 🔥 For local development
        "http://localhost:3001",          // 🔥 Alternative local port
    },
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:          12 * time.Hour,
    }))

    // Health check routes
    router.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message":    "FinBangla Voice Blog API",
            "status":     "success",
            "version":    "1.0.0",
            "mode":       gin.Mode(), // Shows current mode
            "timestamp":  time.Now().Unix(),
        })
    })

    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status":     "healthy",
            "database":   "connected",
            "mode":       gin.Mode(),
            "timestamp":  time.Now().Unix(),
        })
    })

    // Auth routes
    router.GET("/auth/google/login", controllers.GoogleLogin)
    router.GET("/auth/google/callback", controllers.GoogleCallback)

    // Public blog routes
    router.GET("/api/posts", blogController.GetAllPosts)
    router.GET("/api/posts/published", blogController.GetPublishedPosts)
    router.GET("/api/posts/:id", blogController.GetPost)

    // Protected blog routes
    protected := router.Group("/api")
    // TODO: Add authentication middleware when ready
    // protected.Use(middleware.RequireAuth())
    
    protected.POST("/posts", blogController.CreatePost)
    protected.PUT("/posts/:id", blogController.UpdatePost)
    protected.DELETE("/posts/:id", blogController.DeletePost)

    // Comment routes
    router.POST("/api/blogs/:id/comments", commentController.CreateComment)
    router.GET("/api/blogs/:id/comments", commentController.GetCommentsByBlog)
    router.GET("/api/comments/:id", commentController.GetComment)
    router.PUT("/api/comments/:id", commentController.UpdateComment)
    router.DELETE("/api/comments/:id", commentController.DeleteComment)
    router.POST("/api/comments/:id/reply", commentController.CreateReply)
    router.GET("/api/comments/:id/replies", commentController.GetReplies)

    // Get port from environment (Railway sets this automatically)
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("🚀 Server starting on port %s", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatal("❌ Failed to start server:", err)
    }
}

// 🔥 This function checks if all required environment variables exist
func validateEnvironmentVariables() {
    // List of environment variables that MUST exist
    required := []string{
        "GOOGLE_CLIENT_ID",      // For Google OAuth
        "GOOGLE_CLIENT_SECRET",  // For Google OAuth  
        "JWT_SECRET",
        "FRONTEND_URL",           // For generating JWT tokens
    }
    
    // Check each variable one by one
    for _, env := range required {
        // Get the value of environment variable
        value := os.Getenv(env)
        
        // If it's empty or doesn't exist, stop the app
        if value == "" {
            log.Fatalf("❌ Required environment variable %s is not set", env)
        }
    }
    
    log.Println("✅ All required environment variables validated")
}