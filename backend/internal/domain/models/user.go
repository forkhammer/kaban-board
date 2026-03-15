package models

type UserId uint

type User struct {
	IsVisible bool
	IsActive  bool
	Id        UserId
	Name      string
	Username  string
	AvatarUrl string
	Groups    []Group
	Account   *UserAccount
}

type UserAccount struct {
	Id   AccountId
	Name string
	Role AccountRole
}

func (u *User) Validate() error {
	return nil
}
