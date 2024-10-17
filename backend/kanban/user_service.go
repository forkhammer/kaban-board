package kanban

import (
	"fmt"
	"main/config"
	"main/gitlab"
	"main/repository"
	"main/repository/models"
	"strconv"
	"strings"
)

type UserService struct {
	userRepository repository.UserRepositoryInterface `di.inject:"userRepository"`
}

func (s *UserService) CleanUserId(gid string) (uint, error) {
	id, err := strconv.ParseUint(strings.ReplaceAll(gid, "gid://gitlab/User/", ""), 10, 32)

	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func (s *UserService) GetUsers() ([]models.User, error) {
	var users []models.User
	err := s.userRepository.GetUsers(&users)

	if err != nil {
		return users, err
	}

	users = s.cleanUserAvatars(users)
	return users, nil
}

func (s *UserService) GetVisibleUsers() ([]models.User, error) {
	var users []models.User
	err := s.userRepository.GetVisibleUsers(&users)
	return users, err
}

func (s *UserService) saveGitlabUsers(users []gitlab.GitlabUser) ([]models.User, error) {
	var resultUsers = make([]models.User, 0)

	for i := range users {
		u := &users[i]

		userId, err := s.CleanUserId(u.Id)

		if err != nil {
			return resultUsers, err
		}

		var user models.User
		err = s.userRepository.GetOrCreate(&user, models.User{Id: userId}, models.User{
			Id:        userId,
			Name:      u.Name,
			Username:  u.Username,
			AvatarUrl: u.AvatarUrl,
			IsVisible: true,
		})

		if err != nil {
			return resultUsers, err
		}

		resultUsers = append(resultUsers, user)
	}

	return resultUsers, nil
}

func (s *UserService) SetUserVisibility(id int, visibility bool) (*models.User, error) {
	var user models.User
	err := s.userRepository.GetUserBydId(&user, id)

	if err != nil {
		return nil, err
	}

	user.IsVisible = visibility
	err = s.userRepository.SaveUser(user)
	return &user, err
}

func (s *UserService) cleanUserAvatars(users []models.User) []models.User {
	for index := range users {
		user := &(users[index])
		if !strings.HasPrefix(user.AvatarUrl, "https://") {
			user.AvatarUrl = fmt.Sprintf("%s%s", config.Settings.GitlabUrl, user.AvatarUrl)
		}
	}

	return users
}
