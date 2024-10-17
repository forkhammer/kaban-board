package kanban

import (
	"main/gitlab"
	"main/repository"
	"main/repository/models"
	"main/tools"
	"strconv"
	"strings"

	"gorm.io/datatypes"
)

type ProjectService struct {
	projectRepository repository.ProjectRepositoryInterface `di.inject:"projectRepository"`
	userService       UserService
}

func (s *ProjectService) GetProjectById(id uint) (*models.Project, error) {
	var project models.Project
	err := s.projectRepository.GetProjectById(&project, id)
	return &project, err
}

func (s *ProjectService) CleanProjectId(gid string) (uint, error) {
	id, err := strconv.ParseUint(strings.ReplaceAll(gid, "gid://gitlab/Project/", ""), 10, 32)

	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func (s *ProjectService) GetProjects() ([]models.Project, error) {
	var projects []models.Project
	err := s.projectRepository.GetProjects(&projects)
	return projects, err
}

func (s *ProjectService) SaveGitlabProjects(projects []gitlab.GitlabProject) error {
	for _, p := range projects {
		projectId, err := s.CleanProjectId(p.Id)

		if err != nil {
			return err
		}

		userIds := tools.Map(p.GetUserIds(), func(strId string) int64 {
			val, err := s.userService.CleanUserId(strId)
			if err != nil {
				return 0
			}
			return int64(val)
		})

		existProject, err := s.GetProjectById(projectId)
		if err == nil {
			existProject.Name = p.Name
			existProject.Users = datatypes.NewJSONSlice(userIds)
			if err = s.projectRepository.SaveProject(existProject); err != nil {
				return err
			}
		} else {
			project := models.Project{
				Id:        projectId,
				Name:      p.Name,
				IsVisible: true,
				Users:     datatypes.NewJSONSlice(userIds),
			}
			if err = s.projectRepository.CreateProject(&project); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *ProjectService) SetTeam(id uint, teamId *int) (*models.Project, error) {
	project, err := s.GetProjectById(id)

	if err != nil {
		return nil, err
	}

	project.TeamId = teamId
	err = s.projectRepository.SaveProject(project)

	if err != nil {
		return nil, err
	}

	return project, nil
}
