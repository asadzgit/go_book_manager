# go_book_manager
## v1
### Initialize the Go Module:
`go mod init book-management`
### Run the Go Application:
`go run main.go`

### Test the Application:
#### Using curl to Add a Book
`curl -X POST -H "Content-Type: application/json" -d '{"title":"Go Programming","author":"John Doe"}' http://localhost:8080/books
`
#### using curl to see the newly added book
`curl http://localhost:8080/books`

#### Using Postman to Add a Book:

Set the method to POST.

Set the URL to http://localhost:8080/books.

Go to the Body tab, select raw, and choose JSON as the format.

Enter the JSON data:

#### Using Postman to view Books:

make a GET request
http://localhost:8080/books

## v2
handle reading from and writing to a text file instead of array.

## v3
### using database for storage
use `CREATE DATABASE book_management;` in your postgres terminal (`psql`) to create the database.

Then use `go run migrate.go` to run this for creating the books table
