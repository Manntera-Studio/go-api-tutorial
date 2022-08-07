package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo"

	"app/race"
)

type APIMessage struct {
	Message string
}

func main() {
	e := echo.New()
	e.GET("/", dbSelect)
	e.POST("/", dbInsert)
	e.PUT("/:id", dbUpdate)
	e.DELETE("/:id", dbDelete)

	e.Logger.Fatal(e.Start(":8080"))
}

func dbSelect(c echo.Context) error {
	data, err := race.Select(c)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, data)
}

func dbInsert(c echo.Context) error {
	err := race.Insert(c)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	msg := APIMessage{Message: "insert success!"}
	return c.JSON(http.StatusOK, msg)
}

func dbUpdate(c echo.Context) error {
	err := race.Update(c)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	msg := APIMessage{Message: "update success!"}
	return c.JSON(http.StatusOK, msg)
}

func dbDelete(c echo.Context) error {
	err := race.Delete(c)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	msg := APIMessage{Message: "delete success!"}
	return c.JSON(http.StatusOK, msg)
}
