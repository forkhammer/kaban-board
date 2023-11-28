package kanban

import (
	"fmt"
	"main/config"
	"main/gitlab"
	"strconv"
	"strings"
)

type UserService struct {
	userRepository UserRepository
}

func (s *UserService) CleanUserId(gid string) (uint, error) {
	id, err := strconv.ParseUint(strings.ReplaceAll(gid, "gid://gitlab/User/", ""), 10, 32)

	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func (s *UserService) GetUsers() ([]User, error) {
	var repository UserRepository
	users, err := repository.GetUsers()

	if err != nil {
		return users, err
	}

	users = s.cleanUserAvatars(users)
	return users, nil
}

func (s *UserService) GetVisibleUsers() ([]User, error) {
	var repository UserRepository
	return repository.GetVisibleUsers()
}

func (s *UserService) saveGitlabUsers(users []gitlab.GitlabUser) ([]User, error) {
	var resultUsers = make([]User, 0)

	for i := range users {
		u := &users[i]

		userId, err := s.CleanUserId(u.Id)

		if err != nil {
			return resultUsers, err
		}

		user, err := s.userRepository.GetOrCreate(User{Id: userId}, User{
			Id:        userId,
			Name:      u.Name,
			Username:  u.Username,
			AvatarUrl: u.AvatarUrl,
			IsVisible: true,
		})

		if err != nil {
			return resultUsers, err
		}

		resultUsers = append(resultUsers, *user)
	}

	return resultUsers, nil
}

func (s *UserService) SetUserVisibility(id int, visibility bool) (*User, error) {
	user, err := s.userRepository.GetUserBydId(id)

	if err != nil {
		return nil, err
	}

	user.IsVisible = visibility
	err = s.userRepository.SaveUser(user)
	return user, err
}

func (s *UserService) cleanUserAvatars(users []User) []User {
	for index := range users {
		user := &(users[index])
		if !strings.HasPrefix(user.AvatarUrl, "https://") {
			user.AvatarUrl = fmt.Sprintf("%s%s", config.Settings.GitlabUrl, user.AvatarUrl)
		}
	}

	return users
}
