package race

import (
	"strconv"
	"database/sql"
	"log"
	"fmt"
	"os"
	"time"
	"net/url"
	
	"github.com/labstack/echo/v4"
	_ "github.com/go-sql-driver/mysql"
)

func open(path string, count uint) *sql.DB {
	db, err := sql.Open("mysql", path)
	if err != nil {
		log.Fatal("open error:", err)
	}

	if err = db.Ping(); err != nil {
		time.Sleep(time.Second * 2)
		count--
		log.Printf("retry... count:%v", count)
		return open(path, count)
	}

	log.Printf("db connected!")
	return db
}

func connectDB() *sql.DB {
	var path string = fmt.Sprintf("%s:%s@tcp(lesson7-2-db:3306)/%s?charset=utf8&parseTime=true",
		os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_DATABASE"))

	return open(path, 100)
}

func Select(c echo.Context) ([]Race, error) {
	layout := "2006-01-02 15:04:05"
	var id int
	var drid int
	var stime time.Time
	var etime time.Time
	var err error
	if c.QueryParam("id") != "" {
		id, err = strconv.Atoi(c.QueryParam("id"))
	}
	if c.QueryParam("drid") != "" {
		drid, err = strconv.Atoi(c.QueryParam("drid"))
	}
    date := c.QueryParam("date")
	if date == "" && (c.QueryParam("stime") != "" || c.QueryParam("etime") != ""){
		return nil, fmt.Errorf("[ERROR] date is not defined")
	}
	if c.QueryParam("stime") != "" {
		stime, err = time.Parse(layout, fmt.Sprintf("%v %v", date, c.QueryParam("stime")))
	}
	if c.QueryParam("etime") != "" {
		etime, err = time.Parse(layout, fmt.Sprintf("%v %v", date, c.QueryParam("etime")))
	}
	venue := c.QueryParam("venue")
	if err != nil {
		return nil, err
	}
	params := Race{
		Id: id,
		DayRaceId: drid,
		RaceDate: date,
		StartTime: stime,
		EndTime: etime,
		Venue: venue,
	}

	db := connectDB()
	defer db.Close()
	races, err := ReadAll(db, params)
	if err != nil {
		return nil, err
	}
	return races, nil
}

func Insert(c echo.Context) error {
	db := connectDB()
	defer db.Close()

	params, err := c.FormParams()
	if err != nil {
		return err
	}
	layout := "2006-01-02 15:04:05"
	var drid int
	var stime time.Time
	var etime time.Time

	if _, ok := params["drid"]; !ok {
		return customError("drid")
	}
	if _, ok := params["date"]; !ok {
		return customError("date")
	}
	if _, ok := params["stime"]; !ok {
		return customError("stime")
	}
	if _, ok := params["etime"]; !ok {
		return customError("etime")
	}
	if _, ok := params["venue"]; !ok {
		return customError("venue")
	}

	drid, err = strconv.Atoi(params["drid"][0])
	date := params["date"][0]
	stime, err = time.Parse(layout, fmt.Sprintf("%v %v", date, params["stime"][0]))
	etime, err = time.Parse(layout, fmt.Sprintf("%v %v", date, params["etime"][0]))
	venue := params["venue"][0]
	if err != nil {
		return err
	}
	data := Race{
		DayRaceId: drid,
		RaceDate: date,
		StartTime: stime,
		EndTime: etime,
		Venue: venue,
	}
	return insert(db, data)
}

func customError(param string) error {
	return fmt.Errorf("[Error] %v is lacked", param)
}

func Update(c echo.Context) error {
	db := connectDB()
	defer db.Close()

	params, err := c.FormParams()
	if err != nil {
		return err
	}

	var id int
	if c.Param("id") == "" {
		return customError("id")
	}
	id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	layout := "2006-01-02 15:04:05"
	key := updateKey(params)
	if key == "" {
		return fmt.Errorf("[ERROR] invalid update key")
	}
	var drid int
	var date string
	var stime time.Time
	var etime time.Time
	var venue string
	var temperature float64
	if key == "DayRaceId" {
		drid, err = strconv.Atoi(params["drid"][0])
	}
	if key == "RaceDate" {
		date = params["date"][0]
	}
	if key == "StartTime" {
		date = params["date"][0]
		stime, err = time.Parse(layout, fmt.Sprintf("%v %v", date, params["stime"][0]))
	}
	if key == "EndTime" {
		date = params["date"][0]
		etime, err = time.Parse(layout, fmt.Sprintf("%v %v", date, params["etime"][0]))
	}
	if key == "Venue" {
		venue = params["venue"][0]
	}
	if key == "Temperature" {
		temperature, err = strconv.ParseFloat(params["temperature"][0], 64)
	}
	if err != nil {
		return err
	}
	data := Race{
		Id: id,
		DayRaceId: drid,
		RaceDate: date,
		StartTime: stime,
		EndTime: etime,
		Venue: venue,
		Temperature: temperature,
	}

	return update(db, key, data)
}

func updateKey(params url.Values) string {
	if _, ok := params["drid"]; ok {
		return "DayRaceId"
	}
	if _, ok := params["stime"]; ok {
		return "StartTime"
	}
	if _, ok := params["etime"]; ok {
		return "EndTime"
	}
	if _, ok := params["date"]; ok {
		return "RaceDate"
	}
	if _, ok := params["venue"]; ok {
		return "Venue"
	}
	if _, ok := params["temperature"]; ok {
		return "Temperature"
	}
	return ""
}

func Delete(c echo.Context) error {
	db := connectDB()
	defer db.Close()

	i := c.Param("id")
	if i == "" {
		return fmt.Errorf("")
	}

	id, err := strconv.Atoi(i)
	if err != nil {
		return err
	}

	return delete(db, id)
}
