package models

import "time"

type BurndownDataPoint struct {
	Date         time.Time
	TotalScope   uint
	RemainingDev uint
	TotalScopeQA uint
	RemainingQA  uint
}

type BurndownReport struct {
	SprintTitle string
	StartDate   time.Time
	EndDate     time.Time
	DataPoints  []BurndownDataPoint
}
