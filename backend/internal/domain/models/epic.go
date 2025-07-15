package models

type EpicId uint

type Epic struct {
	Id      EpicId
	Title   string
	Project Project
}

func (r *Epic) Validate() error {
	return nil
}
