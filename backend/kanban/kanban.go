package kanban

import (
	"fmt"
	"log"
	"main/cache"
	"main/config"
	"main/gitlab"
	"main/tools"
	"sort"
	"strings"
	"time"

	"github.com/goioc/di"
)

type Kanban struct {
	cache          cache.Cache
	client         *gitlab.GitlabClient
	userService    *UserService
	projectService *ProjectService
	labelService   *LabelService
}

func NewKanban(cache cache.Cache) *Kanban {
	return &Kanban{
		cache:  cache,
		client: gitlab.GitlabClientInstance,
	}
}

func (k *Kanban) PostConstruct() error {
	k.userService = di.GetInstance("userService").(*UserService)
	k.projectService = di.GetInstance("projectService").(*ProjectService)
	k.labelService = di.GetInstance("labelService").(*LabelService)
	return nil
}

func (k *Kanban) GetUsers() ([]KanbanUser, *time.Time, error) {
	users, err := k.userService.GetVisibleUsers()

	if err != nil {
		return make([]KanbanUser, 0), nil, err
	}

	projects, err := k.projectService.GetProjects()

	if err != nil {
		return make([]KanbanUser, 0), nil, err
	}

	issues, err := k.getAllIssues()
	issues = *k.cleanIssues(&issues, &projects)

	if err != nil {
		return make([]KanbanUser, 0), nil, err
	}

	result := make([]KanbanUser, 0)

	for _, user := range users {
		userIssues := tools.Filter(issues, func(issue gitlab.GitlabIssue) bool {
			return tools.IndexOf(issue.Assignees.Nodes, func(a gitlab.GitlabAssignee) bool {
				userId, err := k.userService.CleanUserId(a.UserId)
				return userId == user.Id && err == nil
			}) > -1
		})
		projectIds := tools.Map(userIssues, func(issue gitlab.GitlabIssue) int {
			return issue.ProjectId
		})

		userProjects := tools.Filter(projects, func(project Project) bool {
			hasUser := tools.IndexOf(projectIds, func(id int) bool {
				return uint(id) == project.Id
			}) > -1

			return tools.IndexOf(project.Users, func(id int64) bool {
				return uint(id) == user.Id
			}) > -1 && project.TeamId != nil && hasUser
		})
		teams := tools.Unique(tools.Map(userProjects, func(project Project) uint {
			return uint(*project.TeamId)
		}), func(id uint) uint { return id })

		result = append(result, KanbanUser{
			Id:        user.Id,
			Name:      user.Name,
			Username:  user.Username,
			AvatarUrl: user.AvatarUrl,
			Issues:    userIssues,
			Teams:     teams,
		})
	}
	sort.Slice(result, func(i int, j int) bool {
		return result[i].Name < result[j].Name
	})

	result = *k.cleanUserAvatars(&result)

	var updateTime *time.Time = nil

	if t, foundTime := k.cache.Get("issues_update_time"); foundTime {
		dt := t.(time.Time)
		updateTime = &dt
	}

	return result, updateTime, nil
}

func (k *Kanban) RunUpdater() {
	k.runUpdaterIteration()

	for range time.Tick(time.Minute * time.Duration(config.Settings.GitlabSyncPeriodMin)) {
		k.runUpdaterIteration()
	}
}

func (k *Kanban) runUpdaterIteration() {
	log.Println("Updater tick")

	go k.syncUsers()
	go k.syncProjects()
	go k.syncIssues()
}

func (k *Kanban) getAllIssues() ([]gitlab.GitlabIssue, error) {
	rawResponse, found := k.cache.Get("issues")

	if found {
		return rawResponse.([]gitlab.GitlabIssue), nil
	} else {
		return make([]gitlab.GitlabIssue, 0), nil
	}
}

func (k *Kanban) syncIssues() {
	pageSize := 100
	startCursor := ""
	issues := make([]gitlab.GitlabIssue, 0)

	for {
		response, err := k.client.GetIssuesResponse(pageSize, startCursor)

		if err != nil {
			log.Println(err.Error())
			return
		}

		issues = append(issues, response.Data.Issues.Nodes...)

		if !response.Data.Issues.PageInfo.HasNextPage {
			break
		}

		startCursor = response.Data.Issues.PageInfo.EndCursor
	}

	k.cache.Set("issues", issues, 0)
	k.cache.Set("issues_update_time", time.Now(), 0)

	labels, err := k.extractAllLabels(issues)

	if err != nil {
		log.Println(err.Error())
	}

	_, err = k.labelService.SaveLabels(labels)

	if err != nil {
		log.Println(err.Error())
	}
}

func (k *Kanban) cleanUserAvatars(users *[]KanbanUser) *[]KanbanUser {
	for index := range *users {
		user := &(*users)[index]
		if !strings.HasPrefix(user.AvatarUrl, "https://") {
			user.AvatarUrl = fmt.Sprintf("%s%s", config.Settings.GitlabUrl, user.AvatarUrl)
		}
	}

	return users
}

func (k *Kanban) cleanIssues(issues *[]gitlab.GitlabIssue, projects *[]Project) *[]gitlab.GitlabIssue {
	for index := range *issues {
		issue := &(*issues)[index]
		if issue.Milestone.Id != "" && !strings.HasPrefix(issue.Milestone.WebPath, "https://") {
			issue.Milestone.WebPath = fmt.Sprintf("%s%s", config.Settings.GitlabUrl, issue.Milestone.WebPath)
		}

		projectIndex := tools.IndexOf(*projects, func(p Project) bool {
			return p.Id == uint(issue.ProjectId)
		})

		if projectIndex > -1 {
			project := &(*projects)[projectIndex]
			issue.ProjectName = &project.Name
		}
	}

	return issues
}

func (k *Kanban) syncProjects() {
	projects, err := k.getAllProjects()

	if err != nil {
		log.Println(err.Error())
		return
	}

	if err = k.projectService.SaveGitlabProjects(projects); err != nil {
		log.Println(err.Error())
		return
	}
}

func (k *Kanban) getAllProjects() ([]gitlab.GitlabProject, error) {
	pageSize := 100
	startCursor := ""
	projects := make([]gitlab.GitlabProject, 0)

	for {
		response, err := k.client.GetProjectsResponse(pageSize, startCursor)

		if err != nil {
			return make([]gitlab.GitlabProject, 0), err
		}

		projects = append(projects, response.Data.Projects.Nodes...)

		if !response.Data.Projects.PageInfo.HasNextPage {
			break
		}

		startCursor = response.Data.Projects.PageInfo.EndCursor
	}

	return projects, nil
}

func (k *Kanban) syncUsers() {
	pageSize := 100
	startCursor := ""
	users := make([]gitlab.GitlabUser, 0)

	for {
		response, err := k.client.GetUsersResponse(pageSize, startCursor)

		if err != nil {
			log.Println(err.Error())
			return
		}

		users = append(users, response.Data.Users.Nodes...)

		if !response.Data.Users.PageInfo.HasNextPage {
			break
		}

		startCursor = response.Data.Users.PageInfo.EndCursor
	}

	_, err := k.userService.saveGitlabUsers(users)

	if err != nil {
		log.Println(err.Error())
	}
}

func (k *Kanban) extractAllLabels(issues []gitlab.GitlabIssue) ([]gitlab.GitlabLabel, error) {
	result := make([]gitlab.GitlabLabel, 0)

	for i := range issues {
		issue := &issues[i]
		result = append(result, issue.Labels.Nodes...)
	}

	result = tools.Unique(result, func(issue gitlab.GitlabLabel) string {
		return issue.Id
	})

	return result, nil
}
