package main

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/lib/pq"
)

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

    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatalf("Error opening database: %q", err)
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        log.Fatalf("Cannot connect to the database: %q", err)
    }

    createTable := `
    CREATE TABLE IF NOT EXISTS books (
        id SERIAL PRIMARY KEY,
        title VARCHAR(100),
        author VARCHAR(100)
    );`

    _, err = db.Exec(createTable)
    if err != nil {
        log.Fatalf("Error creating table: %q", err)
    }

    fmt.Println("Migration completed successfully!")
}
