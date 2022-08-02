package main

import (
	"github.com/labstack/echo"
	"app/utils"
)

func main() {
	e := echo.New()
	e.GET("/", utils.Show)

	e.Logger.Fatal(e.Start(":8081"))
}
