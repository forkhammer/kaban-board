package models

type AccountId uint

const (
	AuthProviderPassword AuthProvider = "password"
	AuthProviderGitLab   AuthProvider = "gitlab"
)

type AuthProvider string

type Account struct {
	IsActive     bool
	Id           AccountId
	Username     string
	Password     string
	Name         string
	GitlabID     *uint
	AuthProvider AuthProvider
	AvatarURL    string
}
