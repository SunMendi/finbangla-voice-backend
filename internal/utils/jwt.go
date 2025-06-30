package utils

import (
    "fmt"
    "os"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    Name   string `json:"name"`
    jwt.RegisteredClaims
}

func GenerateJWT(userID uint, email, name string) (string, error) {
    claims := Claims{
        UserID: userID,
        Email:  email,
        Name:   name,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 7)), // 7 days
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "blog-auth-system",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    secretKey := os.Getenv("JWT_SECRET")
    
    if secretKey == "" {
        return "", fmt.Errorf("JWT_SECRET environment variable not set")
    }
    
    return token.SignedString([]byte(secretKey))
}

func ValidateJWT(tokenString string) (*Claims, error) {
    claims := &Claims{}
    secretKey := os.Getenv("JWT_SECRET")
    
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secretKey), nil
    })

    if err != nil {
        return nil, fmt.Errorf("token validation failed: %v", err)
    }

    if !token.Valid {
        return nil, fmt.Errorf("token is not valid")
    }

    return claims, nil
}