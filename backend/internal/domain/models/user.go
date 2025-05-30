package models

type UserId uint

type User struct {
	IsVisible bool
	Id        UserId
	Name      string
	Username  string
	AvatarUrl string
	Groups    []Group
}

func (u *User) Validate() error {
	return nil
}
