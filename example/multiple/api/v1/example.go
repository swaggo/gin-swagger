package v1

import (
	"github.com/gin-gonic/gin"
)

type Book struct {
	ID     int     `json:"id,omitempty"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Year   *uint16 `json:"year"`
}

//
// @Summary Get a list of books in the the store
// @Description get string by ID
// @Accept  json
// @Produce  json
// @Success 200 {array} Book "ok"
// @Router /books [get]
func GetBooks(ctx *gin.Context) {
	ctx.JSON(200, []Book{
		{ID: 1, Title: "Book 1", Author: "Author 1", Year: nil},
		{ID: 2, Title: "Book 2", Author: "Author 2", Year: nil},
	})
}
