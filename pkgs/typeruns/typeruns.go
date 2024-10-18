package typeruns

import (
	"math"
)

type RunInputs struct {
	Value string  `json:"value"`
	Time  float64 `json:"time"`
}

type TypeRun struct {
	Target   string      `json:"target"`
	Html     string      `json:"html"`
	Accuracy float64     `json:"accuracy"`
	Wpm      float64     `json:"wpm"`
	Awpm     float64     `json:"awpm"`
	Time     float64     `json:"time"`
	Inputs   []RunInputs `json:"inputs"`
}

func NewTypeRun() TypeRun {
	return TypeRun{}
}

func Clean(dirty TypeRun) TypeRun {
	clean := dirty

	clean.Accuracy = math.Round(dirty.Accuracy*1000) / 10
	clean.Wpm = math.Round(dirty.Wpm*100) / 100
	clean.Awpm = math.Round(dirty.Awpm*100) / 100
	clean.Time = math.Round(dirty.Time*100) / 100

	return clean
}
