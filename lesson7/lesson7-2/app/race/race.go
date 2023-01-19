package race

import (
	"fmt"
	"time"
	"strings"
	"reflect"
	"database/sql"
)

type Race struct {
	Id int              `json:"id"`
	DayRaceId int       `json:"dayraceid"`
	RaceDate string     `json:"racedate"`
	StartTime time.Time `json:"starttime"`
	EndTime time.Time   `json:"endtime"`
	Temperature float64 `json:"temperature"`
	Venue string        `json:"venuce"`
	Updated time.Time   `json:"updated"`
}

func ReadAll(db *sql.DB, params Race) ([]Race, error) {
	var races []Race
	rows, err := db.Query(makeQuery(params))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		race := Race{}
		err = rows.Scan(&race.Id, &race.DayRaceId, &race.RaceDate, &race.StartTime, &race.EndTime, &race.Temperature, &race.Venue, &race.Updated)
		if err != nil {
			return nil, err
		}
		races = append(races, race)
	}
	rows.Close()

	return races, nil
}

func makeQuery(params Race) string {
	query := "select * from race"

	var cond []string
	if params.Id != 0 {
		cond = append(cond, fmt.Sprintf("id='%v'", params.Id))
	}
	if params.DayRaceId != 0 {
		cond = append(cond, fmt.Sprintf("day_raceid='%v'", params.DayRaceId))
	}
	if params.RaceDate != "" {
		cond = append(cond, fmt.Sprintf("race_date='%v'", params.RaceDate))
	}
	iszero := time.Time{}
	if params.StartTime != iszero {
		cond = append(cond, fmt.Sprintf("start_time='%v'", params.StartTime))
	}
	if params.EndTime != iszero {
		cond = append(cond, fmt.Sprintf("end_time='%v'", params.EndTime))
	}
	if params.Venue != "" {
		cond = append(cond, fmt.Sprintf("venue='%v'", params.Venue))
	}

	if len(cond) > 0 {
		query = fmt.Sprintf("%v where %v;", query, strings.Join(cond, " AND "))
	} else {
		query += ";"
	}

	fmt.Printf("query:%v\n", query)
	return query
}

func insert(db *sql.DB, params Race) error {
	ins, err := db.Prepare("INSERT INTO race (day_raceid, race_date, start_time, end_time, temperature, venue) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer ins.Close()

	_, err = ins.Exec(params.DayRaceId, params.RaceDate, params.StartTime, params.EndTime, params.Temperature, params.Venue)
	if err != nil {
		return err
	}
	return nil
}

func update(db *sql.DB, key string, params Race) error {
	upd, err := db.Prepare(fmt.Sprintf("UPDATE race SET %v = ? WHERE id = ?", key))
	if err != nil {
		return err
	}
	defer upd.Close()

	rv := reflect.ValueOf(params)
	rt := rv.Type()
	name, ok := rt.FieldByName(key)
	if !ok {
		return fmt.Errorf("[ERROR] invalid update key")
	}
	val := rv.FieldByName(name.Name).Interface()
	_, err = upd.Exec(val, params.Id)
	if err != nil {
		return err
	}
	return nil
}

func delete(db *sql.DB, id int) error {
	del, err := db.Prepare("DELETE FROM race WHERE id = ?")
	if err != nil {
		return err
	}
	defer del.Close()

	_, err = del.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
