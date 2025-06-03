package models

import (
	"fmt"
	"strconv"
	"strings"
)

type Quarter struct {
	year   int
	number int
}

func NewQuarter(year int, number int) (*Quarter, error) {
	quarter := &Quarter{
		year:   year,
		number: number,
	}
	if err := quarter.Validate(); err != nil {
		return nil, err
	}
	return quarter, nil
}

func NewQuarterFromId(id string) (*Quarter, error) {
	if len(id) != 7 {
		return nil, fmt.Errorf("invalid id")
	}

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
	return NewQuarter(int(year), int(number))
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
