package models

type TeamId uint

type Team struct {
	Id     TeamId
	Title  string
	Groups []Group
}
