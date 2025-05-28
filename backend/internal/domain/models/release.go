package models

type ReleaseId string
type ReleaseIid string

type Release struct {
	Id        ReleaseId
	Iid       ReleaseIid
	Title     string
	ProjectId ProjectId
}
