package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	bookcrud "github.com/maetad/book-crud"
)

type CreateBookRequest struct {
	Name string `json:"name"`
}

type UpdateBookRequest struct {
	Name string `json:"name"`
}

func main() {
	r := gin.Default()

	b := bookcrud.New("book.sqlite")
	// b.Create()
	// b.Read()
	// b.Update()
	// b.Delete()

	// GET http://localhost:8080/books -> [{"name": "Sapiens"}]
	r.GET("/books", func(c *gin.Context) {
		books := b.Read()

		c.JSON(http.StatusOK, books)
	})

	// POST http://localhost:8080/books {"name": string} -> 201 {"id": uint, "name": string}
	r.POST("/books", func(c *gin.Context) {
		var body CreateBookRequest
		c.ShouldBindJSON(&body)

		if body.Name == "" {
			c.JSON(422, gin.H{
				"message": "name must not empty",
			})
			return
		}

		book := bookcrud.Book{
			Name: body.Name,
		}

		b.Create(&book)

		c.JSON(http.StatusCreated, book)
	})

	// GET http://localhost:8080/books/{id} -> 200 {"id": uint, "name": string}
	r.GET("/books/:id", func(c *gin.Context) {
		var id uint

		cast, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "not found",
			})
			return
		}

		id = uint(cast)

		var body UpdateBookRequest
		c.ShouldBindJSON(&body)

		if body.Name == "" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "name must not empty",
			})
			return
		}

		var book *bookcrud.Book
		books := b.Read()

		// loop to get book.ID == id
		for _, bk := range *books {
			if bk.ID == id {
				book = &bk
				break
			}
		}

		if book.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "not found",
			})
			return
		}

		c.JSON(http.StatusOK, book)
	})

	// PUT http://localhost:8080/books/{id} {"name": string} -> 200 {"id": uint, "name": string}
	r.PUT("/books/:id", func(c *gin.Context) {
		var id uint

		cast, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "not found",
			})
			return
		}

		id = uint(cast)

		var body UpdateBookRequest
		c.ShouldBindJSON(&body)

		if body.Name == "" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "name must not empty",
			})
			return
		}

		var book *bookcrud.Book
		books := b.Read()

		// loop to get book.ID == id
		for _, bk := range *books {
			if bk.ID == id {
				book = &bk
				break
			}
		}

		if book.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "not found",
			})
			return
		}

		book.Name = body.Name
		b.Update(book)

		c.JSON(http.StatusOK, book)
	})

	// DELETE http://localhost:8080/books/{id} -> 204
	r.DELETE("/books/:id", func(c *gin.Context) {
		var id uint

		cast, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "not found",
			})
			return
		}

		id = uint(cast)

		var book *bookcrud.Book
		books := b.Read()

		// loop to get book.ID == id
		for _, bk := range *books {
			if bk.ID == id {
				book = &bk
				break
			}
		}

		if book.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "not found",
			})
			return
		}

		b.Delete(book)

		c.JSON(http.StatusNoContent, nil)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
