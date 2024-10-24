package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/xuri/excelize/v2"
)

var templates *template.Template

// Product represents a single affiliate product
type Products struct {
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

func newProduct(name string, description string, link string, imageurl string) Products {
	// Instantiate the struct and return it
	return Products{
		Name:        name,
		Description: description,
		Link:        link,
		ImageURL:    imageurl,
	}
}

func readProducts() []Products {
	file, err := excelize.OpenFile("static/data/products.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	rows, err := file.GetRows("product")
	if err != nil {
		log.Fatal(err)
	}
	var name, desc, imageURL, link = "", "", "", ""
	var product []Products
	for krow, row := range rows {
		//fmt.Print(krow, "\n")
		if krow != 0 {
			for kcol, col := range row {
				//fmt.Print(kcol, "\n")
				//fmt.Print(col, "\t")
				switch kcol {
				case 1:
					name = col
				case 2:
					desc = col
				case 3:
					imageURL = col
				case 4:
					link = col
				}
			}
			np := newProduct(name, desc, link, imageURL)
			product = append(product, np)
		}
		//fmt.Println()
	}

	return product
}

// homeHandler serves the homepage
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Define some sample products (You can fetch these from an API or database)
	products := readProducts()
	// Render the template and pass the product data
	templates.ExecuteTemplate(w, "index.html", products)
}
