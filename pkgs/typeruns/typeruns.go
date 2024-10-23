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

type ProfileStats struct {
	Awpm     float64 `json:"awpm"`
	Wpm      float64 `json:"wpm"`
	Accuracy float64 `json:"accuracy"`
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
	run := NewTypeRun()
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

func FindByOwner(ownerId int64, db *sql.DB) ([]TypeRun, error) {
	// Query for all runs by owner id
	runQuery := "SELECT * FROM runs WHERE owner_id = ?"
	rows, err := db.Query(runQuery, ownerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Slice to hold all runs for the owner
	var runs []TypeRun

	// Iterate over each row in the result set (each run)
	for rows.Next() {
		run := NewTypeRun()
		err := rows.Scan(&run.Id, &run.OwnerId, &run.Target, &run.Html, &run.Accuracy, &run.Wpm, &run.Awpm, &run.Time)
		if err != nil {
			return nil, err
		}

		// Query for the related RunInputs for each run
		inputsQuery := "SELECT * FROM run_inputs WHERE run_id = ?"
		inputRows, err := db.Query(inputsQuery, run.Id)
		if err != nil {
			return nil, err
		}
		defer inputRows.Close()

		// Populate the Inputs field of the TypeRun
		for inputRows.Next() {
			var input RunInputs
			err := inputRows.Scan(&input.Id, &input.RunId, &input.Value, &input.Time)
			if err != nil {
				return nil, err
			}
			run.Inputs = append(run.Inputs, input)
		}

		// Check for errors in scanning inputs
		if err = inputRows.Err(); err != nil {
			return nil, err
		}

		// Add the run to the list of runs
		runs = append(runs, run.Clean())
	}

	// Check for errors after reading all rows for runs
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return runs, nil
}

func FindBest(count int, db *sql.DB) ([]TypeRun, error) {
	// Query for the best runs ordered by AWPM
	runQuery := "SELECT * FROM runs ORDER BY awpm DESC LIMIT ?"

	rows, err := db.Query(runQuery, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var runs []TypeRun

	// Iterate over each row in the result set (each run)
	for rows.Next() {
		run := NewTypeRun()
		err := rows.Scan(&run.Id, &run.OwnerId, &run.Target, &run.Html, &run.Accuracy, &run.Wpm, &run.Awpm, &run.Time)
		if err != nil {
			return nil, err
		}

		// Query for the related RunInputs for each run
		inputsQuery := "SELECT * FROM run_inputs WHERE run_id = ?"
		inputRows, err := db.Query(inputsQuery, run.Id)
		if err != nil {
			return nil, err
		}
		defer inputRows.Close()

		// Populate the Inputs field of the TypeRun
		for inputRows.Next() {
			var input RunInputs
			err := inputRows.Scan(&input.Id, &input.RunId, &input.Value, &input.Time)
			if err != nil {
				return nil, err
			}
			run.Inputs = append(run.Inputs, input)
		}

		// Check for errors in scanning inputs
		if err = inputRows.Err(); err != nil {
			return nil, err
		}

		// Add the run to the list of best runs
		runs = append(runs, run.Clean())
	}

	// Check for errors after reading all rows for runs
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return runs, nil
}

func CalculateStats(runs []TypeRun) ProfileStats {
	var awpm float64
	var wpm float64
	var accuracy float64
	for _, run := range runs {
		awpm += run.Awpm
		wpm += run.Wpm
		accuracy += run.Accuracy
	}
	count := float64(len(runs))
	awpm /= count
	wpm /= count
	accuracy /= count

	return ProfileStats{
		Awpm:     math.Round(awpm*100) / 100,
		Wpm:      math.Round(wpm*100) / 100,
		Accuracy: math.Round(accuracy*100) / 100,
	}
}
