package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "modernc.org/sqlite" // SQLite driver
)

func main() {
	// Connect to the SQLite database
	db, err := sql.Open("sqlite", "./principles.db")
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	// Start the REPL
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to the Principles Database REPL!")
	fmt.Println("Type `help` for a list of commands.")
	fmt.Println()

	for {
		fmt.Print("> ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		switch {
		case command == "help":
			printHelp()
		case command == "listp":
			listPrinciples(db)
		case command == "listc":
			listCategories(db)
		case strings.HasPrefix(command, "addp"):
			addPrinciple(db, reader)
		case strings.HasPrefix(command, "editp"):
			editPrinciple(db, reader)
		case strings.HasPrefix(command, "linkp"):
			linkPrinciple(db, reader)
		case command == "exit":
			fmt.Println("Exiting the REPL. Goodbye!")
			return
		default:
			fmt.Println("Unknown command. Type `help` for a list of commands.")
		}
	}
}

func printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  list principles            - List all principles")
	fmt.Println("  list categories            - List all categories")
	fmt.Println("  add principle              - Add a new principle")
	fmt.Println("  edit principle [id]        - Edit an existing principle")
	fmt.Println("  link principle [id]        - Link a principle to another principle")
	fmt.Println("  exit                       - Exit the REPL")
}

func listPrinciples(db *sql.DB) {
	rows, err := db.Query(`
		SELECT p.id, p.title, p.description, IFNULL(c.name, 'Uncategorized') 
		FROM principles p
		LEFT JOIN categories c ON p.category_id = c.id
		ORDER BY p.id
	`)
	if err != nil {
		fmt.Println("Error listing principles:", err)
		return
	}
	defer rows.Close()

	fmt.Println("Principles:")
	for rows.Next() {
		var id int
		var title, description, category string
		rows.Scan(&id, &title, &description, &category)
		fmt.Printf("  [%d] %s - %s (Category: %s)\n", id, title, description, category)
	}
}

func listCategories(db *sql.DB) {
	rows, err := db.Query(`SELECT id, name FROM categories ORDER BY id`)
	if err != nil {
		fmt.Println("Error listing categories:", err)
		return
	}
	defer rows.Close()

	fmt.Println("Categories:")
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		fmt.Printf("  [%d] %s\n", id, name)
	}
}

func addPrinciple(db *sql.DB, reader *bufio.Reader) {
	fmt.Print("Enter principle title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Print("Enter principle description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	listCategories(db)
	fmt.Print("Enter category ID (or press Enter to skip): ")
	categoryInput, _ := reader.ReadString('\n')
	categoryInput = strings.TrimSpace(categoryInput)

	var categoryID sql.NullInt64
	if categoryInput != "" {
		id, err := strconv.Atoi(categoryInput)
		if err == nil {
			categoryID.Valid = true
			categoryID.Int64 = int64(id)
		}
	}

	_, err := db.Exec(`INSERT INTO principles (title, description, category_id) VALUES (?, ?, ?)`, title, description, categoryID)
	if err != nil {
		fmt.Println("Error adding principle:", err)
	} else {
		fmt.Println("Principle added successfully!")
	}
}

func editPrinciple(db *sql.DB, reader *bufio.Reader) {
	fmt.Print("Enter the ID of the principle to edit: ")
	idInput, _ := reader.ReadString('\n')
	idInput = strings.TrimSpace(idInput)

	id, err := strconv.Atoi(idInput)
	if err != nil {
		fmt.Println("Invalid ID.")
		return
	}

	var title, description, category string
	err = db.QueryRow(`
		SELECT p.title, p.description, IFNULL(c.name, 'Uncategorized') 
		FROM principles p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = ?
	`, id).Scan(&title, &description, &category)
	if err == sql.ErrNoRows {
		fmt.Println("Principle not found.")
		return
	} else if err != nil {
		fmt.Println("Error retrieving principle:", err)
		return
	}

	fmt.Printf("Current title: %s\n", title)
	fmt.Print("Enter new title (or press Enter to keep): ")
	newTitle, _ := reader.ReadString('\n')
	newTitle = strings.TrimSpace(newTitle)
	if newTitle == "" {
		newTitle = title
	}

	fmt.Printf("Current description: %s\n", description)
	fmt.Print("Enter new description (or press Enter to keep): ")
	newDescription, _ := reader.ReadString('\n')
	newDescription = strings.TrimSpace(newDescription)
	if newDescription == "" {
		newDescription = description
	}

	listCategories(db)
	fmt.Printf("Current category: %s\n", category)
	fmt.Print("Enter new category ID (or press Enter to keep): ")
	newCategoryInput, _ := reader.ReadString('\n')
	newCategoryInput = strings.TrimSpace(newCategoryInput)

	var newCategoryID sql.NullInt64
	if newCategoryInput != "" {
		categoryID, err := strconv.Atoi(newCategoryInput)
		if err == nil {
			newCategoryID.Valid = true
			newCategoryID.Int64 = int64(categoryID)
		}
	} else {
		newCategoryID.Valid = false
	}

	_, err = db.Exec(`UPDATE principles SET title = ?, description = ?, category_id = ? WHERE id = ?`, newTitle, newDescription, newCategoryID, id)
	if err != nil {
		fmt.Println("Error updating principle:", err)
	} else {
		fmt.Println("Principle updated successfully!")
	}
}

func linkPrinciple(db *sql.DB, reader *bufio.Reader) {
	fmt.Print("Enter the ID of the principle to link: ")
	idInput, _ := reader.ReadString('\n')
	idInput = strings.TrimSpace(idInput)

	id, err := strconv.Atoi(idInput)
	if err != nil {
		fmt.Println("Invalid ID.")
		return
	}

	listPrinciples(db)
	fmt.Print("Enter the ID of the related principle: ")
	relatedIDInput, _ := reader.ReadString('\n')
	relatedIDInput = strings.TrimSpace(relatedIDInput)

	relatedID, err := strconv.Atoi(relatedIDInput)
	if err != nil {
		fmt.Println("Invalid related principle ID.")
		return
	}

	_, err = db.Exec(`INSERT INTO principle_links (principle_id, related_id, relation_type) VALUES (?, ?, '')`, id, relatedID)
	if err != nil {
		fmt.Println("Error linking principles:", err)
	} else {
		fmt.Println("Principles linked successfully!")
	}
}
