package db

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Poem struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
}

func ConnectDB() (*sql.DB, error) {
	var dataB *sql.DB
	var err error
	dsn := "root:root@tcp(localhost:3306)/poetry" // Replace with your credentials
	dataB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	if err := dataB.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v\n", err)
	}
	fmt.Println("Connected to MySQL database!")
	return dataB, err
}

func insertRecord() *sql.DB {
	var poem Poem
	dataB, err := ConnectDB()
	if err != nil {
		log.Fatalf("Error encountered: %v", err)
	}
	if err := c.ShouldBindJSON(&poem); err != nil {
		log.Fatalf("Error encountered: %v", err)
	}
	query := "INSERT INTO poems (title, content) VALUES (?, ?)"
	result, err := dataB.Exec(query, poem.Title, poem.Content)
	if err != nil {
		log.Fatalf("Error encountered: %v", err)
	}
	id, _ := result.LastInsertId()
	poem.ID = int(id)
	poem.CreatedAt = time.Now()

	c.JSON(http.StatusCreated, poem)
	return dataB
}
