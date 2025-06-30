package database

import (
    "fmt"
    "log"
    "os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
    //  Railway's DATABASE_URL first
    if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
        var err error
        DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
        if err != nil {
            log.Fatal("Failed to connect to database:", err)
        }
        log.Println("✅ Database connected successfully")
        return
    }
    
    // Fallback to individual environment variables (for local development)
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASS"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )

    // Connect to database
    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Get SQL database for ping
    sqlDB, err := database.DB()
    if err != nil {
        log.Fatal("Failed to get database instance:", err)
    }

    // Test connection
    if err := sqlDB.Ping(); err != nil {
        log.Fatal("Database is not responding:", err)
    }

    // Save connection globally
    DB = database
    log.Println("✅ Database connected successfully!")
}