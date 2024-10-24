package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var templates *template.Template

// Product represents a single affiliate product
type Product struct {
	Name        string
	Description string
	Link        string
	ImageURL    string
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load templates
	templates = template.Must(template.ParseGlob("templates/*.html"))

	// Create a new router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/", homeHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// homeHandler serves the homepage
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Define some sample products (You can fetch these from an API or database)
	products := []Product{
		{Name: "Product 1", Description: "Description for Product 1", Link: "https://example.com/product1", ImageURL: "/static/img/product1.jpg"},
		{Name: "Product 2", Description: "Description for Product 2", Link: "https://example.com/product2", ImageURL: "/static/img/product2.jpg"},
	}

	// Render the template and pass the product data
	templates.ExecuteTemplate(w, "index.html", products)
}
