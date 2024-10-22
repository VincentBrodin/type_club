package typeruns

import (
	"database/sql"
	"math"
)

type RunInputs struct {
	RunId int64   `json:"runId"`
	Value string  `json:"value"`
	Time  float64 `json:"time"`
}

type TypeRun struct {
	Id       int64       `json:"id"`
	Target   string      `json:"target"`
	Html     string      `json:"html"`
	Accuracy float64     `json:"accuracy"`
	Wpm      float64     `json:"wpm"`
	Awpm     float64     `json:"awpm"`
	Time     float64     `json:"time"`
	Inputs   []RunInputs `json:"inputs"`
}

func NewTypeRun() *TypeRun {
	return &TypeRun{
		Id: -1,
	}
}

func (t *TypeRun) Clean() TypeRun {
	clean := *t

	clean.Accuracy = math.Round(t.Accuracy*1000) / 10
	clean.Wpm = math.Round(t.Wpm*100) / 100
	clean.Awpm = math.Round(t.Awpm*100) / 100
	clean.Time = math.Round(t.Time*100) / 100

	return clean
}

func (t *TypeRun) AddToDb(db *sql.DB) error {
	prompt := "INSERT INTO runs (target, html, accuracy, wpm, awpm, time) VALUES (?, ?, ?, ?, ?, ?)"
	statement, err := db.Prepare(prompt)
	if err != nil {
		return err
	}
	defer statement.Close()

	//Execute statement
	result, err := statement.Exec(t.Target, t.Html, t.Accuracy, t.Wpm, t.Awpm, t.Time)
	if err != nil {
		return err
	}

	//Assign the user their id
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	t.Id = id
	return nil
}

func AddToDb(inputs *[]RunInputs, db *sql.DB) error {
	return nil
}
