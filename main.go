package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    _ "github.com/lib/pq"
)

type Book struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
}

var db *sql.DB

const (
    host     = "localhost"
    port     = 5432
    user     = "postgres"
    password = "postgres"
    dbname   = "book_management"
)

func main() {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    var err error
    db, err = sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatalf("Error opening database: %q", err)
    }

    defer db.Close()

    err = db.Ping()
    if err != nil {
        log.Fatalf("Cannot connect to the database: %q", err)
    }

    http.HandleFunc("/books", booksHandler)
    fmt.Println("Server is running on port 8080...")
    http.ListenAndServe(":8080", nil)
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        getBooks(w, r)
    case http.MethodPost:
        createBook(w, r)
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}

func getBooks(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, title, author FROM books")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var books []Book
    for rows.Next() {
        var book Book
        if err := rows.Scan(&book.ID, &book.Title, &book.Author); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        books = append(books, book)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(books)
}

func createBook(w http.ResponseWriter, r *http.Request) {
    var book Book
    if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var lastInsertID int
    err := db.QueryRow("INSERT INTO books(title, author) VALUES($1, $2) RETURNING id", book.Title, book.Author).Scan(&lastInsertID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    book.ID = lastInsertID
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(book)
}
