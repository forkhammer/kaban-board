package models

type EpicId uint

type Epic struct {
	Id         EpicId
	ExternalId IssueExternalId
	Title      string
	Project    Project
}

func (r *Epic) Validate() error {
	return nil
}
