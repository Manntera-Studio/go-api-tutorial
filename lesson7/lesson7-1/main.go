package main

import (
	"net/http"
	"strconv"

	echoSwagger "github.com/swaggo/echo-swagger"
	_ "app/docs"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Student is student strcut
type Student struct {
	ID      int    `json:"id"`
	Grade   int    `json:"grade"`
	Class   string `json:"class"`
	Name    string `json:"name"`
}

// HTTPError is error response
type HTTPError struct {
	Code string `json:"code"`
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample swagger server.
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
func main() {
	e := echo.New()

	initRouting(e)

	e.Logger.Fatal(e.Start(":8080"))
}

func initRouting(e *echo.Echo) {
	// swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/", healthcheck)
	e.GET("/api/v1/classes/:grade/students", getStudents)
}

func healthcheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
}

// getStudents is getting students.
// @Summary get students
// @Description get students in a class
// @Accept  json
// @Produce  json
// @Param grade path int true "Grade"
// @Param class query string false "Class" Enum(A, B)
// @Success 200 {array} main.Student
// @Failure 500 {object} main.HTTPError
// @Router /classes/{grade}/students [get]
func getStudents(c echo.Context) error {
	gradeStr := c.Param("grade")
	grade, err := strconv.Atoi(gradeStr)
	if err != nil {
		return errors.Wrapf(err, "errors when grade convert to int: %s", grade)
	}
	class := c.QueryParam("class")
	students := []*Student{}
	if class == "" || class == "A" {
		students = append(students, &Student{ID: 1, Grade: grade, Class: "A", Name: "Taro"})
		students = append(students, &Student{ID: 2, Grade: grade, Class: "A", Name: "Hanako"})
	}
	if class == "" || class == "B" {
		students = append(students, &Student{ID: 3, Grade: grade, Class: "B", Name: "Jiro"})
		students = append(students, &Student{ID: 4, Grade: grade, Class: "B", Name: "Yuko"})
	}
	return c.JSON(http.StatusOK, students)
}
