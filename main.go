package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type Book struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
}

var books []Book
var idCounter int

func main() {
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
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(books)
}

func createBook(w http.ResponseWriter, r *http.Request) {
    var book Book
    if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    idCounter++
    book.ID = idCounter
    books = append(books, book)
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(book)
}
