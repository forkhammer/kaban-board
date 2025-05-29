package models

type Settings struct {
	TaskTypeLabels []string
}

func (s *Settings) Validate() error {
	return nil
}
