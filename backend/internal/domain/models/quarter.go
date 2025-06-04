package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Quarter struct {
	year   int
	number int
}

func NewQuarter(year int, number int) *Quarter {
	quarter := &Quarter{
		year:   year,
		number: number,
	}
	return quarter
}

func NewQuarterFromId(id string) (*Quarter, error) {
	split := strings.Split(id, "-")
	if len(split) != 2 {
		return nil, fmt.Errorf("invalid id separator")
	}

	year, err := strconv.ParseInt(split[0], 10, 32)
	if err != nil {
		return nil, err
	}
	number, err := strconv.ParseInt(split[1], 10, 32)
	if err != nil {
		return nil, err
	}
	quarter := NewQuarter(int(year), int(number))
	return quarter, quarter.Validate()
}

func NewQuarterFromDate(dt time.Time) *Quarter {
	return NewQuarter(dt.Year(), int((dt.Month()-1)/3+1))
}

func (q *Quarter) Validate() error {
	if q.year <= 0 {
		return fmt.Errorf("year must be greater than 0")
	}
	if q.number <= 0 {
		return fmt.Errorf("number must be greater than 0")
	}
	if q.number > 4 {
		return fmt.Errorf("number must be less than 4")
	}
	return nil
}

func (q *Quarter) GetId() string {
	return fmt.Sprintf("%d-%d", q.year, q.number)
}

func (q *Quarter) GetTitle() string {
	return fmt.Sprintf("%d квартал %d года", q.number, q.year)
}

func (q *Quarter) StartDate() time.Time {
	return time.Date(q.year, time.Month((q.number-1)*3+1), 1, 0, 0, 0, 0, time.Local)
}

func (q *Quarter) EndDate() time.Time {
	return time.Date(q.year, time.Month((q.number-1)*3+4), 1, 0, 0, 0, 0, time.Local).Add(time.Duration(-1) * time.Minute)
}
