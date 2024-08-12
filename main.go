package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// in Go | type | response (Fiber handle this)
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book

func main() {
	app := fiber.New()

	// Initialize in-memory data
	books = append(books, Book{ID: 1, Title: "1984", Author: "George Orwell"})
	books = append(books, Book{ID: 2, Title: "The Great Gatsby", Author: "F. Scott Fitzgerald"})

	app.Get("/books", func(c *fiber.Ctx) error {
		return c.JSON(books)
	})

	// Use the logging middleware
	app.Use(loggingMiddleware)
	// NOTE: Middleware will not works at GET /books

	app.Get("book/:id", getBook)
	app.Post("/book", createBook)
	app.Put("/book/:id", updateBook)

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := os.Getenv("SECRET")
	if secretKey == "" {
		secretKey = "defaultSecret" // Default value if not specified
	}

	app.Get("/secret", func(c *fiber.Ctx) error {
		return c.SendString(secretKey)
	})

	// NOTE: nodemon does not reload the env, have to restart the server

	app.Listen(":8080")
}

func getBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	for _, book := range books {
		if book.ID == id {
			return c.JSON(book)
		}
	}

	return c.SendStatus(fiber.StatusNotFound)
}

func createBook(c *fiber.Ctx) error {
	book := new(Book)

	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	book.ID = len(books) + 1
	books = append(books, *book)

	return c.JSON(book)
}

func updateBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	bookUpdate := new(Book)
	if err := c.BodyParser(bookUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for i, book := range books {
		if book.ID == id {
			book.Title = bookUpdate.Title
			book.Author = bookUpdate.Author
			books[i] = book
			return c.JSON(book)
		}
	}

	return c.SendStatus(fiber.StatusNotFound)
}

func loggingMiddleware(c *fiber.Ctx) error {
	// Start timer
	start := time.Now()

	// Process request
	err := c.Next()

	// Calculate processing time
	duration := time.Since(start)

	// Log the information
	fmt.Printf("Request URL: %s - Method: %s - Duration: %s\n", c.OriginalURL(), c.Method(), duration)

	return err
}
