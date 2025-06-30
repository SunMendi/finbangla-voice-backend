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
        
        // Get client IP
        clientIP := c.ClientIP()
        method := c.Request.Method
        path := c.Request.URL.Path
        
        // Log request start
        log.Printf("🔍 %s %s from %s", method, path, clientIP)
        
        // Process request
        c.Next()
        
        // Calculate duration
        duration := time.Since(start)
        statusCode := c.Writer.Status()
        
        // Log request completion
        log.Printf("✅ %s %s | %d | %v | %s", 
            method, path, statusCode, duration, clientIP)
    }
}