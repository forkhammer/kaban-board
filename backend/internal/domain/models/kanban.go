package models

import (
	"main/pkg/utils"
	"time"
)

type KanbanUser struct {
	Issues []Issue
	Teams  []Team
	User   *User
}

type Board struct {
	Users      []KanbanUser
	UpdateTime *time.Time
}

type Kanban struct {
	Users  []User
	Issues []Issue
}

func NewKanban(users []User, issues []Issue) (*Kanban, error) {
	kanban := &Kanban{
		Users:  users,
		Issues: issues,
	}

	if err := kanban.Validate(); err != nil {
		return nil, err
	}
	return kanban, nil
}

func (k *Kanban) Validate() error {
	return nil
}

func (k *Kanban) GetBoard() (*Board, error) {
	kanbanUsers := make([]KanbanUser, 0, len(k.Users))
	for _, user := range k.Users {

		issues := utils.Filter(k.Issues, func(issue Issue) bool {
			return utils.IndexOf(issue.Assignees, func(assignee User) bool {
				return assignee.Id == user.Id
			}) > -1
		})

		projects := utils.Unique(
			utils.Map(issues, func(issue Issue) Project {
				return issue.Project
			}),
			func(project Project) uint {
				return uint(project.Id)
			},
		)

		teams := utils.Unique(
			utils.Map(
				utils.Filter(projects, func(project Project) bool {
					return project.Team != nil
				}),
				func(project Project) Team {
					return *project.Team
				},
			),
			func(team Team) uint {
				return uint(team.Id)
			},
		)

		kanbanUser := KanbanUser{
			User:   &user,
			Teams:  teams,
			Issues: []Issue{},
		}
		kanbanUsers = append(kanbanUsers, kanbanUser)
	}
	now := time.Now()
	return &Board{
		Users:      kanbanUsers,
		UpdateTime: &now,
	}, nil
}
