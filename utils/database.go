package utils

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func OpenDB() *sql.DB {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	// Create the data source name (DSN)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	db, _ := sql.Open("postgres", dsn)
	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db.PingContext(ctx)
	fmt.Println("success")
	return db
}
