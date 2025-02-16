package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "modernc.org/sqlite" // SQLite driver
)

type Category struct {
	ID   int
	Name string
}

// Principle represents a principle with title, description, and category
type Principle struct {
	ID               int
	Title            string
	Description      template.HTML
	Category         string
	CategoryId       int
	PrevID           int         // Previous principle ID
	NextID           int         // Next principle ID
	LinkedPrinciples []Principle // Linked principles
}

func handlePrinciple(w http.ResponseWriter, r *http.Request) {
	// Load the template
	tmpl := template.Must(template.ParseFiles("templates/principle.html"))
	// Connect to SQLite database
	db, err := sql.Open("sqlite", "./principles.db")
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get the "id" query parameter
	idParam := r.URL.Query().Get("id")
	var currentID int
	if idParam == "" {
		// Default to the first principle if no "id" is provided
		row := db.QueryRow("SELECT id FROM principles ORDER BY id LIMIT 1")
		err = row.Scan(&currentID)
		if err != nil {
			http.Error(w, "No principles found", http.StatusNotFound)
			return
		}
	} else {
		// Parse the "id" parameter to an integer
		currentID, err = strconv.Atoi(idParam)
		if err != nil || currentID < 1 {
			http.Error(w, "Invalid principle ID", http.StatusBadRequest)
			return
		}
	}

	// Fetch the current principle
	var principle Principle
	row := db.QueryRow(`
			SELECT p.id, p.title, p.description, IFNULL(c.name, 'Uncategorized'), c.id AS category
			FROM principles p
			LEFT JOIN categories c ON p.category_id = c.id
			WHERE p.id = ?
		`, currentID)
	err = row.Scan(&principle.ID, &principle.Title, &principle.Description, &principle.Category, &principle.CategoryId)
	if err == sql.ErrNoRows {
		http.Error(w, "Principle not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch linked principles
	rows, err := db.Query(`
		SELECT p.id, p.title, p.description, IFNULL(c.name, 'Uncategorized'), c.id AS category
		FROM principles p
		LEFT JOIN categories c ON p.category_id = c.id
		INNER JOIN principle_links pl ON p.id = pl.related_id
		WHERE pl.principle_id = ?
		UNION
		SELECT p.id, p.title, p.description, IFNULL(c.name, 'Uncategorized'), c.id AS category
		FROM principles p
		LEFT JOIN categories c ON p.category_id = c.id
		INNER JOIN principle_links pl ON p.id = pl.principle_id
		WHERE pl.related_id = ?
	`, principle.ID, principle.ID)
	if err != nil {
		http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var linkedPrinciples []Principle
	for rows.Next() {
		var linkedPrinciple Principle
		if err := rows.Scan(&linkedPrinciple.ID, &linkedPrinciple.Title, &linkedPrinciple.Description, &linkedPrinciple.Category, &linkedPrinciple.CategoryId); err != nil {
			http.Error(w, "Database scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		linkedPrinciples = append(linkedPrinciples, linkedPrinciple)
	}
	principle.LinkedPrinciples = linkedPrinciples

	// Calculate Previous and Next IDs
	// Fetch the previous principle ID
	row = db.QueryRow("SELECT id FROM principles WHERE id < ? ORDER BY id DESC LIMIT 1", principle.ID)
	row.Scan(&principle.PrevID) // If no result, PrevID will remain 0

	// Fetch the next principle ID
	row = db.QueryRow("SELECT id FROM principles WHERE id > ? ORDER BY id ASC LIMIT 1", principle.ID)
	row.Scan(&principle.NextID) // If no result, NextID will remain 0

	// Render the template
	tmpl.Execute(w, principle)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// Fetch categories from the database
	db, err := sql.Open("sqlite", "./principles.db")
	rows, err := db.Query("SELECT id, name FROM categories")
	if err != nil {
		http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			http.Error(w, "Database scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		http.Error(w, "Database rows error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Pass categories to the template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, categories)
}

func handleCategory(w http.ResponseWriter, r *http.Request) {
	// Connect to SQLite database
	db, err := sql.Open("sqlite", "./principles.db")
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get the "id" query parameter
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "Category ID is required", http.StatusBadRequest)
		return
	}

	// Parse the "id" parameter to an integer
	categoryID, err := strconv.Atoi(idParam)
	if err != nil || categoryID < 1 {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// Fetch the current category
	var category Category
	row := db.QueryRow("SELECT id, name FROM categories WHERE id = ?", categoryID)
	err = row.Scan(&category.ID, &category.Name)
	if err == sql.ErrNoRows {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch principles for the current category
	rows, err := db.Query("SELECT id, title FROM principles WHERE category_id = ?", categoryID)
	if err != nil {
		http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var principles []Principle
	for rows.Next() {
		var principle Principle
		if err := rows.Scan(&principle.ID, &principle.Title); err != nil {
			http.Error(w, "Database scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		principles = append(principles, principle)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		http.Error(w, "Database rows error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch all categories for navigation
	allCategoriesRows, err := db.Query("SELECT id, name FROM categories")
	if err != nil {
		http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer allCategoriesRows.Close()

	var allCategories []Category
	for allCategoriesRows.Next() {
		var cat Category
		if err := allCategoriesRows.Scan(&cat.ID, &cat.Name); err != nil {
			http.Error(w, "Database scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		allCategories = append(allCategories, cat)
	}

	// Check for errors from iterating over rows
	if err := allCategoriesRows.Err(); err != nil {
		http.Error(w, "Database rows error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Pass data to the template
	data := struct {
		CurrentCategory Category
		Categories      []Category
		Principles      []Principle
	}{
		CurrentCategory: category,
		Categories:      allCategories,
		Principles:      principles,
	}

	tmpl, err := template.ParseFiles("templates/category.html")
	if err != nil {
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}

func handleCss(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/style.css")
}

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/style.css", handleCss)
	http.HandleFunc("/principle", handlePrinciple)
	http.HandleFunc("/category", handleCategory)

	// Start the HTTP server
	log.Println("Server running at http://localhost:5001/")
	log.Fatal(http.ListenAndServe(":5001", nil))
}
