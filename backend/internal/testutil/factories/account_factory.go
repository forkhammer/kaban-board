package factories

import (
	"main/internal/app/interfaces"
	"main/internal/domain/models"
	"main/internal/domain/repo"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/goioc/di"
)

type AccountFactory struct {
	repo   repo.AccountRepo
	hasher interfaces.PasswordServiceInterface
}

func NewAccountFactory() *AccountFactory {
	return &AccountFactory{
		repo:   di.GetInstance("AccountRepository").(repo.AccountRepo),
		hasher: di.GetInstance("PasswordService").(interfaces.PasswordServiceInterface),
	}
}

func (f *AccountFactory) Build(data map[string]any) *models.Account {
	account := &models.Account{}

	if val, ok := data["username"].(string); ok && val != "" {
		account.Username = val
	} else {
		account.Username = gofakeit.Username()
	}

	if val, ok := data["name"].(string); ok && val != "" {
		account.Name = val
	} else {
		account.Name = gofakeit.Name()
	}

	if val, ok := data["password"].(string); ok && val != "" {
		hashed, _ := f.hasher.HashPassword(val)
		account.Password = hashed
	} else {
		plain := gofakeit.Password(true, true, true, true, false, 16)
		hashed, _ := f.hasher.HashPassword(plain)
		account.Password = hashed
	}

	if val, ok := data["is_active"].(bool); ok {
		account.IsActive = val
	} else {
		account.IsActive = true
	}

	if val, ok := data["role"].(string); ok && val != "" {
		account.Role = models.AccountRole(val)
	} else {
		account.Role = models.AccountRoleViewer
	}

	if val, ok := data["avatar_url"].(string); ok && val != "" {
		account.AvatarURL = val
	} else {
		account.AvatarURL = gofakeit.URL()
	}

	if val, ok := data["gitlab_id"].(*models.UserId); ok {
		account.GitlabID = val
	} else if val, ok := data["gitlab_id"].(uint); ok {
		id := models.UserId(val)
		account.GitlabID = &id
	}

	salt, _ := f.hasher.GenerateSalt()
	account.JwtSalt = salt
	account.AuthProvider = models.AuthProviderPassword

	return account
}

func (f *AccountFactory) Create(data map[string]any) (*models.Account, error) {
	account := f.Build(data)
	return f.repo.Create(account)
}
