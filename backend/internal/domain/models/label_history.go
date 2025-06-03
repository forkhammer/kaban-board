package models

import "time"

type LabelHistory struct {
	Labels    []LabelId
	CreatedAt time.Time
}

func (l *LabelHistory) Validate() error {
	return nil
}
