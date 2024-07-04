package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strconv"
    "strings"
)

type Book struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
}

var books []Book
var idCounter int
const fileName = "books.txt"

func main() {
    loadBooks()

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
    saveBooks()
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(book)
}

func loadBooks() {
    file, err := os.Open(fileName)
    if err != nil {
        if os.IsNotExist(err) {
            return
        }
        fmt.Println("Error opening file:", err)
        return
    }
    defer file.Close()

    data, err := ioutil.ReadAll(file)
    if err != nil {
        fmt.Println("Error reading file:", err)
        return
    }

    lines := strings.Split(string(data), "\n")
    for _, line := range lines {
        if line == "" {
            continue
        }
        parts := strings.Split(line, ",")
        id, _ := strconv.Atoi(parts[0])
        book := Book{
            ID:     id,
            Title:  parts[1],
            Author: parts[2],
        }
        books = append(books, book)
        if id > idCounter {
            idCounter = id
        }
    }
}

func saveBooks() {
    var data strings.Builder
    for _, book := range books {
        line := fmt.Sprintf("%d,%s,%s\n", book.ID, book.Title, book.Author)
        data.WriteString(line)
    }

    err := ioutil.WriteFile(fileName, []byte(data.String()), 0644)
    if err != nil {
        fmt.Println("Error writing to file:", err)
    }
}
