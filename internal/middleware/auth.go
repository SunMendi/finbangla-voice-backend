package middleware

import (
    "auth2_google/internal/utils"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "error":   "Authorization header required",
            })
            c.Abort()
            return
        }

        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

        _, err := utils.ValidateJWT(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "error":   "Invalid or expired token",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}