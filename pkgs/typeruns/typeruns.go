package typeruns

import (
	"database/sql"
	"math"
)

type RunInputs struct {
	Id    int64   `json:"id"`
	RunId int64   `json:"runId"`
	Value string  `json:"value"`
	Time  float64 `json:"time"`
}

type TypeRun struct {
	Id       int64       `json:"id"`
	OwnerId  int64       `json:"ownerId"`
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
	prompt := "INSERT INTO runs (owner_id, target, html, accuracy, wpm, awpm, time) VALUES (?, ?, ?, ?, ?, ?, ?)"
	statement, err := db.Prepare(prompt)
	if err != nil {
		return err
	}
	defer statement.Close()

	//Execute statement
	result, err := statement.Exec(t.OwnerId, t.Target, t.Html, t.Accuracy, t.Wpm, t.Awpm, t.Time)
	if err != nil {
		return err
	}

	//Assign the user their id
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	t.Id = id
	for i := range t.Inputs {
		t.Inputs[i].RunId = id
	}
	AddInputsToDb(t.Inputs, db)
	return nil
}

func AddInputsToDb(inputs []RunInputs, db *sql.DB) error {
	for i := range inputs {
		prompt := "INSERT INTO run_inputs (run_id, value, time) VALUES (?, ?, ?)"
		statement, err := db.Prepare(prompt)
		if err != nil {
			return err
		}
		defer statement.Close()

		//Execute statement
		result, err := statement.Exec(inputs[i].RunId, inputs[i].Value, inputs[i].Time)
		if err != nil {
			return err
		}

		//Assign the user their id
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		inputs[i].Id = id

	}
	return nil
}

func FindById(id int64, db *sql.DB) (*TypeRun, error) {
	// Query for the TypeRun by id
	runQuery := "SELECT * FROM runs WHERE id = ?"
	run := &TypeRun{}
	err := db.QueryRow(runQuery, id).Scan(&run.Id, &run.OwnerId, &run.Target, &run.Html, &run.Accuracy, &run.Wpm, &run.Awpm, &run.Time)
	if err != nil {
		return nil, err
	}

	// Query for the related RunInputs
	inputsQuery := "SELECT * FROM run_inputs WHERE run_id = ?"
	rows, err := db.Query(inputsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Populate the Inputs field of the TypeRun
	for rows.Next() {
		var input RunInputs
		err := rows.Scan(&input.Id, &input.RunId, &input.Value, &input.Time)
		if err != nil {
			return nil, err
		}
		run.Inputs = append(run.Inputs, input)
	}

	// Check for any errors after reading all rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return run, nil
}
