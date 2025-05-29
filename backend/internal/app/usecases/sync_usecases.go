package usecases

import (
	"fmt"
	"main/internal/app/interfaces"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type SyncUseCases struct {
	gitlab   interfaces.TaskTracker `di.inject:"gitlab"`
	userRepo repo.UserRepo          `di.inject:"UserRepository"`
}

func (uc *SyncUseCases) Sync() error {
	err := uc.SyncProjects()
	if err != nil {
		return fmt.Errorf("Error syncing projects: %w", err)
	}

	err = uc.SyncUsers()
	if err != nil {
		return fmt.Errorf("Error syncing users: %w", err)
	}

	err = uc.SyncIssues()
	if err != nil {
		return fmt.Errorf("Error syncing issues: %w", err)
	}
	return nil
}

func (uc *SyncUseCases) SyncProjects() error {
	return nil
}

func (uc *SyncUseCases) SyncUsers() error {
	users, err := uc.gitlab.GetUsers()
	if err != nil {
		return err
	}

	for _, user := range users {
		existUser, err := uc.userRepo.Get(domain.UserId(user.Id))
		if err == nil {
			existUser.Name = user.Name
			existUser.Username = user.Username
			existUser.AvatarUrl = user.AvatarUrl

			if err := existUser.Validate(); err != nil {
				return err
			}

			_, err := uc.userRepo.Update(existUser)
			if err != nil {
				return err
			}
		} else {
			_, err := uc.userRepo.Create(&user)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (uc *SyncUseCases) SyncIssues() error {
	return nil
}
