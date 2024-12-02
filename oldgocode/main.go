package main

import (
	db "affiliate-website/connect"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var templates *template.Template

const file string = "digest.db"

type Category struct {
	Name        string
	Description string
	Link        string
	ImageURL    string
}
type Data struct {
	Category   []map[string]interface{}
	Affiliates []Category
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

// GetUsers retrieves all user records from the database
func GetCategories(db *sql.DB, sqlStr string) ([]map[string]interface{}, error) {
	rows, err := db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		var description string
		var imageURL string
		var link string
		if err := rows.Scan(&id, &name, &description, &imageURL, &link); err != nil {
			return nil, err
		}
		category := map[string]interface{}{
			"id":          id,
			"name":        name,
			"description": description,
			"imageURL":    imageURL,
			"link":        link,
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func newCategory(name string, description string, link string, imageurl string) Category {
	// Instantiate the struct and return it
	return Category{
		Name:        name,
		Description: description,
		Link:        link,
		ImageURL:    imageurl,
	}
}

func readCategory() []map[string]interface{} {
	dbVar, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbVar.Close()
	var sql string = "SELECT * FROM category"
	categories, err := GetCategories(dbVar, sql)
	if err != nil {
		log.Fatal("Failed to get categories:", err)
	}

	return categories
}

// homeHandler serves the homepage
func homeHandler(w http.ResponseWriter, r *http.Request) {
	var data Data
	data.Category = readCategory()
	data.Affiliates = nil
	tmpl, err := template.New("").ParseFiles("templates/index.html", "templates/base.html")
	if err != nil {
		log.Fatal(err)
	}
	// Render the template and pass the category data
	tmpl.ExecuteTemplate(w, "base", data)
}
