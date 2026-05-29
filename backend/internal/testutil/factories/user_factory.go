package factories

import (
	"main/internal/domain/models"
	"main/internal/domain/repo"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/goioc/di"
)

type UserFactory struct {
	repo repo.UserRepo
}

func NewUserFactory() *UserFactory {
	return &UserFactory{
		repo: di.GetInstance("UserRepository").(repo.UserRepo),
	}
}

func (f *UserFactory) Build(data map[string]any) *models.User {
	user := &models.User{}

	if val, ok := data["name"].(string); ok && val != "" {
		user.Name = val
	} else {
		user.Name = gofakeit.Name()
	}

	if val, ok := data["username"].(string); ok && val != "" {
		user.Username = val
	} else {
		user.Username = gofakeit.Username()
	}

	if val, ok := data["avatar_url"].(string); ok && val != "" {
		user.AvatarUrl = val
	} else {
		user.AvatarUrl = gofakeit.URL()
	}

	if val, ok := data["is_visible"].(bool); ok {
		user.IsVisible = val
	} else {
		user.IsVisible = true
	}

	if val, ok := data["is_active"].(bool); ok {
		user.IsActive = val
	} else {
		user.IsActive = true
	}

	if val, ok := data["groups"].([]models.Group); ok {
		user.Groups = val
	}

	return user
}

func (f *UserFactory) Create(data map[string]any) (*models.User, error) {
	user := f.Build(data)
	return f.repo.Create(user)
}
