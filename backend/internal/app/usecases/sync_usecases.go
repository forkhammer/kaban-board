package usecases

import (
	"errors"
	"fmt"
	"main/config"
	"main/internal/app/interfaces"
	domain_pkg "main/internal/domain"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type SyncUseCases struct {
	gitlab          interfaces.TaskTracker              `di.inject:"gitlab"`
	config          *config.Config                      `di.inject:"config"`
	userRepo        repo.UserRepo                       `di.inject:"UserRepository"`
	accountRepo     repo.AccountRepo                    `di.inject:"AccountRepository"`
	passwordService interfaces.PasswordServiceInterface `di.inject:"PasswordService"`
	projectRepo     repo.ProjectRepo                    `di.inject:"ProjectRepository"`
	issueRepo       repo.IssueRepo                      `di.inject:"IssueRepository"`
	labelRepo       repo.LabelRepo                      `di.inject:"LabelRepository"`
	releaseRepo     repo.ReleaseRepo                    `di.inject:"ReleaseRepository"`
	epicRepo        repo.EpicRepo                       `di.inject:"EpicRepository"`
}

func (uc *SyncUseCases) Sync() error {
	err := uc.SyncProjects()
	if err != nil {
		return fmt.Errorf("Error syncing projects: %w", err)
	}

	err = uc.SyncReleases()
	if err != nil {
		return fmt.Errorf("Error syncing releases: %w", err)
	}

	err = uc.SyncUsers()
	if err != nil {
		return fmt.Errorf("Error syncing users: %w", err)
	}

	err = uc.SyncLabels()
	if err != nil {
		return fmt.Errorf("Error syncing labels: %w", err)
	}

	err = uc.SyncIssues()
	if err != nil {
		return fmt.Errorf("Error syncing issues: %w", err)
	}
	return nil
}

func (uc *SyncUseCases) SyncProjects() error {
	projects, err := uc.gitlab.GetProjects()
	if err != nil {
		return err
	}

	for _, project := range projects {
		existProject, err := uc.projectRepo.Get(project.Id)
		if err == nil {
			existProject.Name = project.Name
			existProject.Users = project.Users

			if err := existProject.Validate(); err != nil {
				return err
			}

			_, err := uc.projectRepo.Update(existProject)
			if err != nil {
				return err
			}
		} else {
			_, err := uc.projectRepo.Create(&project)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (uc *SyncUseCases) SyncReleases() error {
	releases, err := uc.gitlab.GetReleases()
	if err != nil {
		return err
	}

	for _, release := range releases {
		existRelease, err := uc.releaseRepo.Get(release.Id)
		if err == nil {
			existRelease.Title = release.Title
			existRelease.Iid = release.Iid
			existRelease.Project = release.Project
			existRelease.WebPath = release.WebPath

			if err := existRelease.Validate(); err != nil {
				return err
			}

			_, err := uc.releaseRepo.Update(existRelease)
			if err != nil {
				return err
			}
		} else {
			_, err := uc.releaseRepo.Create(&release)
			if err != nil {
				return err
			}
		}
	}

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
			existUser.IsActive = user.IsActive

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

		if err := uc.syncAccountActivity(user); err != nil {
			return err
		}
	}

	return nil
}

func (uc *SyncUseCases) syncAccountActivity(user domain.User) error {
	account, err := uc.accountRepo.GetByGitlabID(uint(user.Id))
	if err != nil {
		return nil
	}

	if account.IsActive == user.IsActive {
		return nil
	}

	wasActive := account.IsActive
	account.IsActive = user.IsActive

	if wasActive && !user.IsActive {
		newSalt, err := uc.passwordService.GenerateSalt()
		if err != nil {
			return fmt.Errorf("failed to generate jwt salt: %w", err)
		}
		account.JwtSalt = newSalt
	}

	_, err = uc.accountRepo.Update(account)
	return err
}

func (uc *SyncUseCases) SyncIssues() error {
	issues, err := uc.gitlab.GetIssues()
	if err != nil {
		return err
	}

	for _, issue := range issues {
		existIssue, err := uc.issueRepo.GetByExternalId(issue.ExternalId)

		if err != nil {
			if _, ok := errors.AsType[*domain_pkg.NotFoundError](err); !ok {
				return err
			}
		}

		if existIssue != nil {
			existIssue.Iid = issue.Iid
			existIssue.ExternalId = issue.ExternalId
			existIssue.Title = issue.Title
			existIssue.IssueType = issue.IssueType
			existIssue.Assignees = issue.Assignees
			existIssue.WebUrl = issue.WebUrl
			existIssue.Labels = issue.Labels
			existIssue.Project = issue.Project
			existIssue.Release = issue.Release
			existIssue.TaskType = issue.TaskType
			existIssue.EstimateDev = issue.EstimateDev
			existIssue.EstimateQA = issue.EstimateQA

			if err := existIssue.Validate(); err != nil {
				return err
			}

			_, err := uc.issueRepo.Update(existIssue)
			if err != nil {
				return err
			}
		} else {
			_, err := uc.issueRepo.Create(&issue)
			if err != nil {
				return err
			}
		}
	}

	return uc.syncEpicsFromIssues(issues)
}

func (uc *SyncUseCases) syncEpicsFromIssues(issues []domain.Issue) error {
	epicLabelText := uc.config.GitlabEpicLabel
	for _, issue := range issues {
		isEpic := false
		for _, label := range issue.Labels {
			if label.Name == epicLabelText {
				isEpic = true
				break
			}
		}
		if !isEpic {
			continue
		}

		externalId := issue.ExternalId
		existEpic, err := uc.epicRepo.GetByExternalId(externalId)
		if err != nil {
			if _, ok := errors.AsType[*domain_pkg.NotFoundError](err); !ok {
				return err
			}
		}

		if existEpic != nil {
			existEpic.Iid = issue.Iid
			existEpic.Title = issue.Title
			existEpic.Project = issue.Project
			if _, err := uc.epicRepo.Update(existEpic); err != nil {
				return err
			}
		} else {
			newEpic := &domain.Epic{
				ExternalId: externalId,
				Iid:        issue.Iid,
				Title:      issue.Title,
				Project:    issue.Project,
			}
			if _, err := uc.epicRepo.Create(newEpic); err != nil {
				return err
			}
		}
	}
	return nil
}

func (uc *SyncUseCases) SyncLabels() error {
	labels, err := uc.gitlab.GetLabels()
	if err != nil {
		return err
	}

	for _, label := range labels {
		existLabel, err := uc.labelRepo.Get(label.Id)
		if err == nil {
			existLabel.Name = label.Name
			existLabel.Color = label.Color
			existLabel.TextColor = label.TextColor

			if err := existLabel.Validate(); err != nil {
				return err
			}

			_, err := uc.labelRepo.Update(existLabel)
			if err != nil {
				return err
			}
		} else {
			_, err := uc.labelRepo.Create(&label)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
