package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// serialize struct
type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgeral", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func getBooks(c *gin.Context) {
	// Gin context is everything about the http request and allows you to give a response
	// books will automatically be serialized into json from struct
	c.IndentedJSON(http.StatusOK, books)
}

func bookById(c *gin.Context) {
	// Param means Path Parameter i.e /books/2
	id := c.Param("id")
	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	c.IndentedJSON(http.StatusOK, book)

}

func checkoutBook(c *gin.Context) {
	// Gets query parameter ?id=2
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id parameter"})
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not found."})
		return
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing ID Parameter"})
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not found."})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)

}

func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			// Return a pointer of the book so we can actually modify the fields of the struct from
			// a different function
			return &books[i], nil
		}
	}
	return nil, errors.New("book not found")
}

func createBook(c *gin.Context) {
	var newBook book

	// BindJson binds the data part of the data payload in the request
	// to the book struct
	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}
func main() {
	router := gin.Default()
	router.GET("/books", getBooks)

	// To get path paramater us :
	router.GET("/books/:id", bookById)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:8080")

}
