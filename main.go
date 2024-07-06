package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "github.com/gorilla/mux"
    "html/template"
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

    router := mux.NewRouter()
    fmt.Println("Server is running on port 8080...")
    
    // we’re using the `template.ParseFiles` function to parse the HTML template. 
    // The `tmpl.Execute` function is used to render the template and write the result to the `http.ResponseWriter` interface. 
    // In this case, we’re passing the name “Mehul” to the template. 
    router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl, err := template.ParseFiles("templates/home.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        err = tmpl.Execute(w, "Asad")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    })
    
    router.HandleFunc("/books", booksHandler)
    http.ListenAndServe(":8080", router)
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

    query_params := r.URL.Query() /* getting query params from url */
    fmt.Println("%s", query_params)
    // fmt.Fprintf(w, "%s", query_params) /* sending simple plain text */
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
