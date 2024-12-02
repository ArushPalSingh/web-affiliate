package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var dataB *sql.DB

func generatePoem(prompt string) (string, error) {
	apiURL := "https://api-inference.huggingface.co/models/bigscience/bloom"
	apiToken := os.Getenv("HUGGINGFACE_API_TOKEN")

	// Create payload
	payload := map[string]string{"inputs": prompt}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	// Extract text
	if generatedText, ok := result["generated_text"].(string); ok {
		return generatedText, nil
	}

	return "", fmt.Errorf("no text generated")
}

func insertPoem(prompt string) {
	var err error
	text, err := generatePoem(prompt)
	if err != nil {
		log.Fatalf("Error generating text: %v", err)
	}
	db.insertRecord()
	defer dataB.Close()

}

func main() {
	prompt := "Write a short poem about the sunrise."
	insertPoem(prompt)
	fmt.Println("Generated Text:", text)
}
