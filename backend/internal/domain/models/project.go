package models

type ProjectId uint

type Project struct {
	IsVisible bool
	Id        ProjectId
	Name      string
	TeamId    *TeamId
	Users     []UserId
}
