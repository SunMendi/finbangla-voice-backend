package middleware

import (
    "log"
    "time"
    
    "github.com/gin-gonic/gin"
)

// BasicLogger - Simple middleware for IP and request time
func BasicLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Record start time
        start := time.Now()
        
        // Get client IP and request details
        clientIP := c.ClientIP()
        method := c.Request.Method
        path := c.Request.URL.Path
        origin := c.GetHeader("Origin")
        userAgent := c.GetHeader("User-Agent")
        
        // Log request start with Origin
        if origin != "" {
            log.Printf("🔍 %s %s from %s | Origin: %s", 
                method, path, clientIP, origin)
        } else {
            log.Printf("🔍 %s %s from %s | No Origin header", 
                method, path, clientIP)
        }
        
        // Log User-Agent for additional debugging
        if userAgent != "" {
            log.Printf("🔍 User-Agent: %s", userAgent)
        }
        
        // Process request
        c.Next()
        
        // Calculate duration
        duration := time.Since(start)
        statusCode := c.Writer.Status()
        
        // Log request completion with Origin
        if origin != "" {
            log.Printf("✅ %s %s | %d | %v | %s | Origin: %s", 
                method, path, statusCode, duration, clientIP, origin)
        } else {
            log.Printf("✅ %s %s | %d | %v | %s", 
                method, path, statusCode, duration, clientIP)
        }
    }
}