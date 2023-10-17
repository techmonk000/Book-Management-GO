package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/yugabyte/ysql-go"
	_ "github.com/yugabyte/ysql-go"
)

const (
	hostname = "localhost"
	port     = 5433
	database = "bookstore"
	user     = "your_db_user"
	password = "your_db_password"
)

var db *sql.DB

func main() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", hostname, port, user, password, database)
	var err error
	db, err = sql.Open("ysql", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/books", listBooks)
	http.HandleFunc("/add-book", addBook)

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func listBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, author, year FROM books")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		if err != nil {
			log.Fatal(err)
		}
		books = append(books, book)
	}

	fmt.Fprint(w, "List of Books:\n")
	for _, book := range books {
		fmt.Fprintf(w, "Title: %s, Author: %s, Year: %d\n", book.Title, book.Author, book.Year)
	}
}

func addBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		author := r.FormValue("author")
		year := r.FormValue("year")

		_, err := db.Exec("INSERT INTO books (title, author, year) VALUES ($1, $2, $3)", title, author, year)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/books", http.StatusSeeOther)
		return
	}

	fmt.Fprint(w, `
		<html>
		<body>
		<h1>Add a Book</h1>
		<form method="post" action="/add-book">
			Title: <input type="text" name="title"><br>
			Author: <input type="text" name="author"><br>
			Year: <input type="text" name="year"><br>
			<input type="submit" value="Add Book">
		</form>
		</body>
		</html>
	`)
}

type Book struct {
	ID     int
	Title  string
	Author string
	Year   int
}
