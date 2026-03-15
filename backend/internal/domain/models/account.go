package models

type AccountId uint

type AuthProvider string

const (
	AuthProviderPassword AuthProvider = "password"
	AuthProviderGitLab   AuthProvider = "gitlab"
)

type AccountRole string

const (
	AccountRoleAdmin    AccountRole = "admin"
	AccountRoleEmployee AccountRole = "employee"
	AccountRoleViewer   AccountRole = "viewer"
)

type Account struct {
	IsActive     bool
	Id           AccountId
	Username     string
	Password     string
	Name         string
	GitlabID     *UserId
	AuthProvider AuthProvider
	AvatarURL    string
	Role         AccountRole
	JwtSalt      string
}
