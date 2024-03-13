package main

import (
    "database/sql"
    "log"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    _ "github.com/mattn/go-sqlite3"
)

type Book struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
    Author string `json:"author"`
}

var db *sql.DB

func main() {
    var err error
    db, err = sql.Open("sqlite3", "./books.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    createTable()

    router := gin.Default()

    // Routes
    router.GET("/books", getBooks)
    router.GET("/books/:id", getBook)
    router.POST("/books", createBook)
    router.PUT("/books/:id", updateBook)
    router.DELETE("/books/:id", deleteBook)

    // Start server
    if err := router.Run(":8080"); err != nil {
        log.Fatal(err)
    }
}

func createTable() {
    query := `
    CREATE TABLE IF NOT EXISTS books (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        author TEXT NOT NULL
    );
    `
    if _, err := db.Exec(query); err != nil {
        log.Fatal(err)
    }
}

func getBooks(c *gin.Context) {
    var books []Book
    rows, err := db.Query("SELECT id, title, author FROM books")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    for rows.Next() {
        var book Book
        if err := rows.Scan(&book.ID, &book.Title, &book.Author); err != nil {
            log.Println(err)
            continue
        }
        books = append(books, book)
    }

    c.JSON(http.StatusOK, books)
}

func getBook(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    var book Book
    err := db.QueryRow("SELECT id, title, author FROM books WHERE id = ?", id).Scan(&book.ID, &book.Title, &book.Author)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
        return
    }

    c.JSON(http.StatusOK, book)
}

func createBook(c *gin.Context) {
    var book Book
    if err := c.BindJSON(&book); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    result, err := db.Exec("INSERT INTO books (title, author) VALUES (?, ?)", book.Title, book.Author)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    id, _ := result.LastInsertId()
    book.ID = int(id)
    c.JSON(http.StatusCreated, book)
}

func updateBook(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    var book Book
    if err := c.BindJSON(&book); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    _, err := db.Exec("UPDATE books SET title = ?, author = ? WHERE id = ?", book.Title, book.Author, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    book.ID = id
    c.JSON(http.StatusOK, book)
}

func deleteBook(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    _, err := db.Exec("DELETE FROM books WHERE id = ?", id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}
