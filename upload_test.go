package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"log"
	"os"
	"testing"
)

func SaveImage(db *sql.DB, name string, img []byte) error {
	query := `INSERT INTO image(name, image) VALUES($1, $2)`
	_, err := db.Exec(query, name, img)
	if err != nil {
		return err
	}
	return nil
}

func ViewImage(db *sql.DB, id int64) (string, []byte, error) {
	var imgName string
	var imgData []byte
	query := `SELECT name, image FROM image WHERE id = $1;`
	err := db.QueryRow(query, id).Scan(&imgName, &imgData)
	if err != nil {
		return "", nil, err
	}
	return imgName, imgData, nil
}

func TestViewImage(t *testing.T) {
	// Open connection to database
	connStr := "host=localhost port=5432 user=spotify password=pa55word dbname=spotify sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// View image
	name, img, err := ViewImage(db, 4)
	fmt.Println(img)
	err = os.WriteFile(fmt.Sprintf("C:/Users/adit/Documents/View/%s.jpg", name), img, 0644)
	if err != nil {
		log.Fatalf("Failed to save image to file: %v", err)
	}

	fmt.Sprintf("Hexadecimal data successfully converted to binary and saved as '%s.jpg'.", name)
}

func TestUploadImage(t *testing.T) {
	// Open connection to database
	connStr := "host=localhost port=5432 user=spotify password=pa55word dbname=spotify sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// Get image and read
	imgPath := "C:/Users/adit/Downloads/coba.jpeg"
	imgFile, err := os.Open(imgPath)
	if err != nil {
		log.Println(err)
	}
	imgData, err := io.ReadAll(imgFile)
	if err != nil {
		log.Println(err)
	}

	// Insert image to database
	err = SaveImage(db, "try1", imgData)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Saved!")

}
