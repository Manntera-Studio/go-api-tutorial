package main

import (
	"log"
	"net/http"
	"fmt"

	echoSwagger "github.com/swaggo/echo-swagger"
	_ "app/docs"
	"github.com/labstack/echo/v4"

	"app/race"
)

type APIMessage struct {
	Message string `json:"message"`
}

// @title Racing Database API
// @version 1.0
// @description Racing Database
// @license.name MIT License
// @license.url https://opensource.org/licenses/mit-license.php
// @host localhost:8081
// @BasePath /
func main() {
	e := echo.New()

	initRouting(e)

	e.Logger.Fatal(e.Start(":8081"))
}

func initRouting(e *echo.Echo) {
	// swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/", dbSelect)
	e.POST("/", dbInsert)
	e.PUT("/:id", dbUpdate)
	e.DELETE("/:id", dbDelete)

}

// dbSelect is select data.
// @Summary select data
// @Description select data from database
// @Acctept json
// @Produce json
// @Param id query int false "Id"
// @Param drid query int false "DayRaceId"
// @Param date query string false "RaceDate" Format(date)
// @Param stime query string false "StartTime"
// @Param etime query string false "EndTime"
// @Param venue query string false "Venue"
// @Success 200 {array} race.Race
// @Failure 400 {object} main.APIMessage
// @Failure 500 {object} main.APIMessage
// @Router / [get]
func dbSelect(c echo.Context) error {
	data, err := race.Select(c)
	if err != nil {
		log.Printf("%v", err)
		msg := APIMessage{Message: fmt.Sprintf("%s", err)}
		return c.JSON(http.StatusBadRequest, msg)
	}
	return c.JSON(http.StatusOK, data)
}

// dbInsert is insert data.
// @Summary insert data
// @Description insert data from database
// @Acctept json
// @Produce json
// @Param drid query int true "DayRaceId"
// @Param date query string true "RaceDate" Format(date)
// @Param stime query string true "StartTime"
// @Param venue query string true "Venue"
// @Param etime query string true "EndTime"
// @Success 200 {object} main.APIMessage
// @Failure 400 {object} main.APIMessage
// @Failure 500 {object} main.APIMessage
// @Router / [post]
func dbInsert(c echo.Context) error {
	err := race.Insert(c)
	if err != nil {
		log.Printf("%v", err)
		msg := APIMessage{Message: fmt.Sprintf("%s", err)}
		return c.JSON(http.StatusBadRequest, msg)
	}
	msg := APIMessage{Message: "insert success!"}
	return c.JSON(http.StatusOK, msg)
}

// dbInsert is update data.
// @Summary update data
// @Description update data from database
// @Acctept json
// @Produce json
// @Param id path int true "Id"
// @Param drid query int false "DayRaceId"
// @Param date query string false "RaceDate"
// @Param stime query string false "StartTime"
// @Param etime query string false "EndTime"
// @Param temperature query number false "Temperature"
// @Param venue query string false "Venue"
// @Success 200 {object} main.APIMessage
// @Failure 400 {object} main.APIMessage
// @Failure 500 {object} main.APIMessage
// @Router /{id} [put]
func dbUpdate(c echo.Context) error {
	err := race.Update(c)
	if err != nil {
		log.Printf("%v", err)
		msg := APIMessage{Message: fmt.Sprintf("%s", err)}
		return c.JSON(http.StatusBadRequest, msg)
	}
	msg := APIMessage{Message: "update success!"}
	return c.JSON(http.StatusOK, msg)
}

// dbInsert is update data.
// @Summary update data
// @Description update data from database
// @Acctept json
// @Produce json
// @Param id path int true "Id"
// @Success 200 {object} main.APIMessage
// @Failure 400 {object} main.APIMessage
// @Failure 500 {object} main.APIMessage
// @Router /{id} [delete]
func dbDelete(c echo.Context) error {
	err := race.Delete(c)
	if err != nil {
		log.Printf("%v", err)
		msg := APIMessage{Message: fmt.Sprintf("%s", err)}
		return c.JSON(http.StatusBadRequest, msg)
	}
	msg := APIMessage{Message: "delete success!"}
	return c.JSON(http.StatusOK, msg)
}
