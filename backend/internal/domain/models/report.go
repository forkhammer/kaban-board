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

type BurnupDataPoint struct {
	Date         time.Time
	ScopeDev     uint
	CompletedDev uint
	ScopeQA      uint
	CompletedQA  uint
}

type BurnupReport struct {
	SprintTitle string
	StartDate   time.Time
	EndDate     time.Time
	DataPoints  []BurnupDataPoint
}
