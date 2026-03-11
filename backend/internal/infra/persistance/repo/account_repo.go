package repo

import (
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"

	"gorm.io/gorm"
)

type AccountRepository struct {
	conn interfaces.ConnectionInterface `di.inject:"db"`
}

func NewAccountRepository(conn interfaces.ConnectionInterface) *AccountRepository {
	return &AccountRepository{
		conn: conn,
	}
}

func (r *AccountRepository) Get(id domain.AccountId) (*domain.Account, error) {
	account := &models.Account{}
	if err := r.conn.GetEngine().Where("accounts.id = ?", id).First(account).Error; err != nil {
		return nil, err
	}
	return r.toDomainAccount(account), nil
}

func (r *AccountRepository) GetByUsername(username string) (*domain.Account, error) {
	account := &models.Account{}
	if err := r.conn.GetEngine().Where("accounts.username = ?", username).First(account).Error; err != nil {
		return nil, err
	}
	return r.toDomainAccount(account), nil
}

func (r *AccountRepository) GetByGitlabID(gitlabID uint) (*domain.Account, error) {
	account := &models.Account{}
	if err := r.conn.GetEngine().Where("accounts.gitlab_id = ?", gitlabID).First(account).Error; err != nil {
		return nil, err
	}
	return r.toDomainAccount(account), nil
}

func (r *AccountRepository) List(spec repo.QuerySpec) ([]domain.Account, error) {
	accounts := make([]models.Account, 0)
	query := r.conn.GetEngine().Model(&models.Account{})

	if spec != nil {
		if result, err := spec.Apply(query); err != nil {
			return nil, err
		} else {
			query = result.(*gorm.DB)
		}
	}

	if err := query.Find(&accounts).Error; err != nil {
		return nil, err
	}

	domainAccounts := make([]domain.Account, len(accounts))
	for i, account := range accounts {
		domainAccounts[i] = *r.toDomainAccount(&account)
	}

	return domainAccounts, nil
}

func (r *AccountRepository) Count(spec repo.QuerySpec) (int, error) {
	var count int64
	query := r.conn.GetEngine().Model(&models.Account{})

	if spec != nil {
		if result, err := spec.Apply(query); err != nil {
			return 0, err
		} else {
			query = result.(*gorm.DB)
		}
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *AccountRepository) Create(account *domain.Account) (*domain.Account, error) {
	model := r.toAccount(account)
	err := r.conn.GetEngine().Create(model).Error
	if err != nil {
		return nil, err
	}
	return r.Get(domain.AccountId(model.Id))
}

func (r *AccountRepository) Update(account *domain.Account) (*domain.Account, error) {
	model := r.toAccount(account)
	err := r.conn.GetEngine().Save(model).Error
	if err != nil {
		return nil, err
	}
	return r.Get(domain.AccountId(model.Id))
}

func (r *AccountRepository) Delete(id domain.AccountId) error {
	return r.conn.GetEngine().Where("id = ?", id).Delete(&models.Account{}).Error
}

func (r *AccountRepository) toDomainAccount(account *models.Account) *domain.Account {
	var gitlabId *domain.UserId
	if account.GitlabID != nil {
		id := domain.UserId(*account.GitlabID)
		gitlabId = &id
	}

	return &domain.Account{
		Id:           domain.AccountId(account.Id),
		Username:     account.Username,
		Password:     account.Password,
		Name:         account.Name,
		IsActive:     account.IsActive,
		GitlabID:     gitlabId,
		AuthProvider: domain.AuthProvider(account.AuthProvider),
		AvatarURL:    account.AvatarURL,
		Role:         domain.AccountRole(account.Role),
	}
}

func (r *AccountRepository) toAccount(account *domain.Account) *models.Account {
	var gitlabId *uint
	if account.GitlabID != nil {
		id := uint(*account.GitlabID)
		gitlabId = &id
	}

	return &models.Account{
		Id:           uint(account.Id),
		Username:     account.Username,
		Password:     account.Password,
		Name:         account.Name,
		IsActive:     account.IsActive,
		GitlabID:     gitlabId,
		AuthProvider: string(account.AuthProvider),
		AvatarURL:    account.AvatarURL,
		Role:         string(account.Role),
	}
}
