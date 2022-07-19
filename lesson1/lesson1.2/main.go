package main

import (
	"net/http"

	"github.com/labstack/echo"
)

type User struct {
	Name string
	Email string
}

func main() {
	e := echo.New()
	e.GET("/user", show)

	e.Logger.Fatal(e.Start(":5000"))
}

func show(c echo.Context) error {
    name := c.QueryParam("name")
    email := c.QueryParam("email")

	u := new(User)
	u.Name = name
	u.Email = email

	return c.JSON(http.StatusOK, u)
}
