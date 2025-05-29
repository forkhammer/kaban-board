package models

type ProjectId uint

type Project struct {
	IsVisible bool
	Id        ProjectId
	Name      string
	Team      *Team
	Users     []UserId
}

func (p *Project) Validate() error {
	return nil
}
