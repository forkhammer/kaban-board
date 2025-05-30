package models

type ReleaseId string
type ReleaseIid string

type Release struct {
	Id      ReleaseId
	Iid     ReleaseIid
	Title   string
	Project Project
}

func (r *Release) Validate() error {
	return nil
}
