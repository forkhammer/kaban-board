package utils

import (
	"fmt"
	"math"
)

func IntToUintPtr(p *int) (*uint, error) {
	if p == nil {
		return nil, nil
	}
	if *p < 0 {
		return nil, fmt.Errorf("Cannot convert negative int to uint")
	}
	u := uint(*p)
	return &u, nil
}

func UintToIntPtr(p *uint) (*int, error) {
	if p == nil {
		return nil, nil
	}
	if *p > math.MaxInt {
		return nil, fmt.Errorf("Cannot convert uint larger than int max value")
	}
	i := int(*p)
	return &i, nil
}
