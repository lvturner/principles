package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "modernc.org/sqlite" // SQLite driver
)

// Principle represents a principle with title, description, and category
type Principle struct {
	Title       string
	Description string
	Category    string
}

func main() {
	// Load the HTML template
	tmpl := template.Must(template.ParseFiles("templates/principle.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Connect to the SQLite database
		db, err := sql.Open("sqlite", "./principles.db")
		if err != nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// Get the "id" query parameter from the request
		idParam := r.URL.Query().Get("id")
		var principle Principle

		if idParam == "" {
			// If no ID is provided, show the first principle by default
			row := db.QueryRow(`
                SELECT p.title, p.description, IFNULL(c.name, 'Uncategorized') AS category
                FROM principles p
                LEFT JOIN categories c ON p.category_id = c.id
                ORDER BY p.id LIMIT 1
            `)
			err = row.Scan(&principle.Title, &principle.Description, &principle.Category)
		} else {
			// Convert the "id" parameter to an integer
			id, err := strconv.Atoi(idParam)
			if err != nil || id < 1 {
				http.Error(w, "Invalid principle ID", http.StatusBadRequest)
				return
			}

			// Query the database for the specified principle ID
			row := db.QueryRow(`
                SELECT p.title, p.description, IFNULL(c.name, 'Uncategorized') AS category
                FROM principles p
                LEFT JOIN categories c ON p.category_id = c.id
                WHERE p.id = ?
            `, id)
			err = row.Scan(&principle.Title, &principle.Description, &principle.Category)
		}

		// Handle errors from the database query
		if err == sql.ErrNoRows {
			http.Error(w, "Principle not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Database query error", http.StatusInternalServerError)
			return
		}

		// Render the template with the principle
		tmpl.Execute(w, principle)
	})

	// Start the HTTP server
	fmt.Println("Server running at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
