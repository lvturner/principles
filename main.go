package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "modernc.org/sqlite" // SQLite driver
)

// Principle represents a principle with title, description, and category
type Principle struct {
	ID          int
	Title       string
	Description string
	Category    string
	PrevID      int // Previous principle ID
	NextID      int // Next principle ID
}

func handlePrinciple(w http.ResponseWriter, r *http.Request) {
	// Load the template
	tmpl := template.Must(template.ParseFiles("templates/principle.html"))
		// Connect to SQLite database
		db, err := sql.Open("sqlite", "./principles.db")
		if (err != nil) {
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
			SELECT p.id, p.title, p.description, IFNULL(c.name, 'Uncategorized') AS category
			FROM principles p
			LEFT JOIN categories c ON p.category_id = c.id
			WHERE p.id = ?
		`, currentID)
		err = row.Scan(&principle.ID, &principle.Title, &principle.Description, &principle.Category)
		if err == sql.ErrNoRows {
			http.Error(w, "Principle not found", http.StatusNotFound)
			return
		} else if err != nil {
      http.Error(w, "Database query error: " +err.Error(), http.StatusInternalServerError)
			return
		}

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
	http.ServeFile(w, r, "templates/index.html")
}

func handleCss(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/style.css")
}

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/style.css", handleCss)
	http.HandleFunc("/principle", handlePrinciple) 

	// Start the HTTP server
	log.Println("Server running at http://localhost:5001/")
	log.Fatal(http.ListenAndServe(":5002", nil))
}
