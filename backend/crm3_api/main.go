package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func hello(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

type createTodoRequest struct {
	Title string `json:"title"`
}

type createTodoResponse struct {
	ID int `json:"id"`
}

func createTodo(c echo.Context) error {
	req := new(createTodoRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, &createTodoResponse{ID: 1})
}

func main() {
	e := echo.New()

	e.GET("/", hello)
	e.POST("/todo", createTodo)

	e.Start(":1323")
}