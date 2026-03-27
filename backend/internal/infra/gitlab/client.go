package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"main/config"
	"main/internal/app/interfaces"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/pkg/utils"

	"github.com/goioc/di"
)

type GitlabClient struct {
	interfaces.TaskTracker
	apiUrl      string
	token       string
	config      *config.Config    `di.inject:"config"`
	userRepo    repo.UserRepo     `di.inject:"UserRepository"`
	projectRepo repo.ProjectRepo  `di.inject:"ProjectRepository"`
	labelRepo   repo.LabelRepo    `di.inject:"LabelRepository"`
	settingRepo repo.SettingsRepo `di.inject:"SettingsRepository"`
	releaseRepo repo.ReleaseRepo  `di.inject:"ReleaseRepository"`
}

func NewGitlabClient(apiUrl string, token string) *GitlabClient {
	return &GitlabClient{
		apiUrl: apiUrl,
		token:  token,
	}
}

func (client *GitlabClient) PostConstruct() error {
	client.config = di.GetInstance("config").(*config.Config)
	client.userRepo = di.GetInstance("UserRepository").(repo.UserRepo)
	client.projectRepo = di.GetInstance("ProjectRepository").(repo.ProjectRepo)
	client.labelRepo = di.GetInstance("LabelRepository").(repo.LabelRepo)
	client.settingRepo = di.GetInstance("SettingsRepository").(repo.SettingsRepo)
	client.releaseRepo = di.GetInstance("ReleaseRepository").(repo.ReleaseRepo)
	return nil
}

func (client *GitlabClient) GetProjects() ([]domain.Project, error) {
	pageSize := 100
	startCursor := ""
	projects := make([]GitlabProject, 0)

	for {
		response, err := client.GetProjectsResponse(pageSize, startCursor)

		if err != nil {
			return []domain.Project{}, err
		}

		projects = append(projects, response.Data.Projects.Nodes...)

		if !response.Data.Projects.PageInfo.HasNextPage {
			break
		}

		startCursor = response.Data.Projects.PageInfo.EndCursor
	}

	result := make([]domain.Project, 0)
	for _, project := range projects {
		domainProject, err := client.toDomainProject(&project)
		if err != nil {
			return []domain.Project{}, err
		}
		result = append(result, *domainProject)
	}

	return result, nil
}

func (client *GitlabClient) GetUsers() ([]domain.User, error) {
	pageSize := 100
	startCursor := ""
	users := make([]GitlabUser, 0)

	for {
		response, err := client.GetUsersResponse(pageSize, startCursor)

		if err != nil {
			log.Println(err.Error())
			return []domain.User{}, err
		}

		users = append(users, response.Data.Users.Nodes...)

		if !response.Data.Users.PageInfo.HasNextPage {
			break
		}

		startCursor = response.Data.Users.PageInfo.EndCursor
	}

	domainUsers := make([]domain.User, 0)

	for _, user := range users {
		domainUser, err := client.toDomainUser(&user)

		if err != nil {
			return []domain.User{}, err
		}

		domainUsers = append(domainUsers, *domainUser)
	}

	return domainUsers, nil
}

func (client *GitlabClient) GetIssues() ([]domain.Issue, error) {
	pageSize := 100
	startCursor := ""
	issues := make([]GitlabIssue, 0)

	for {
		response, err := client.GetIssuesResponse(pageSize, startCursor)

		if err != nil {
			return []domain.Issue{}, err
		}

		issues = append(issues, response.Data.Issues.Nodes...)

		if !response.Data.Issues.PageInfo.HasNextPage {
			break
		}

		startCursor = response.Data.Issues.PageInfo.EndCursor
	}

	users, err := client.userRepo.List(nil)
	if err != nil {
		return []domain.Issue{}, err
	}

	projects, err := client.projectRepo.List(nil)
	if err != nil {
		return []domain.Issue{}, err
	}

	labels, err := client.labelRepo.List(nil)
	if err != nil {
		return []domain.Issue{}, err
	}

	settings, err := client.settingRepo.Get()
	if err != nil {
		return []domain.Issue{}, err
	}

	releases, err := client.releaseRepo.List(nil)
	if err != nil {
		return []domain.Issue{}, err
	}

	result := make([]domain.Issue, 0)
	for _, issue := range issues {
		domainIssue, err := client.toDomainIssue(
			&issue,
			users,
			projects,
			labels,
			releases,
			settings,
		)
		if err != nil {
			return []domain.Issue{}, err
		}
		result = append(result, *domainIssue)
	}

	return result, nil
}

func (client *GitlabClient) GetLabels() ([]domain.Label, error) {
	pageSize := 100
	startCursor := ""
	issues := make([]GitlabIssue, 0)

	for {
		response, err := client.GetIssuesResponse(pageSize, startCursor)

		if err != nil {
			return []domain.Label{}, err
		}

		issues = append(issues, response.Data.Issues.Nodes...)

		if !response.Data.Issues.PageInfo.HasNextPage {
			break
		}

		startCursor = response.Data.Issues.PageInfo.EndCursor
	}

	labels := make([]GitlabLabel, 0)
	for _, issue := range issues {
		labels = append(labels, issue.Labels.Nodes...)
	}
	labels = utils.Unique(labels, func(issue GitlabLabel) string {
		return issue.Id
	})

	domainLabels := utils.Map(labels, func(label GitlabLabel) domain.Label {
		return *client.toDomainLabel(&label)
	})
	return domainLabels, nil
}

func (client *GitlabClient) GetReleases() ([]domain.Release, error) {
	pageSize := 100
	startCursor := ""
	projects := make([]GitlabProject, 0)

	for {
		response, err := client.GetProjectsResponse(pageSize, startCursor)

		if err != nil {
			return []domain.Release{}, err
		}

		projects = append(projects, response.Data.Projects.Nodes...)

		if !response.Data.Projects.PageInfo.HasNextPage {
			break
		}

		startCursor = response.Data.Projects.PageInfo.EndCursor
	}

	milesones := make([]MilestoneIndex, 0)
	for _, project := range projects {
		for _, node := range project.Milesones.Nodes {
			milesones = append(milesones, MilestoneIndex{
				Milestone: &node,
				Project:   &project,
			})
		}
	}
	milesones = utils.Unique(milesones, func(milestone MilestoneIndex) string {
		return milestone.Milestone.Id
	})

	result := make([]domain.Release, 0)
	for _, milesone := range milesones {
		domainProject, err := client.toDomainProject(milesone.Project)
		if err != nil {
			return []domain.Release{}, err
		}

		domainRelease, err := client.toDomainRelease(milesone.Milestone, domainProject)
		if err != nil {
			return []domain.Release{}, err
		}
		result = append(result, *domainRelease)
	}

	return result, nil
}

func (client *GitlabClient) GetUsersResponse(pageSize int, startCursor string) (*GitlabUsersResponse, error) {
	data, err := client.graphQLRequest(client.getUsersQuery(pageSize, startCursor))

	if err != nil {
		return nil, err
	}

	var response GitlabUsersResponse

	if err = json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (client *GitlabClient) GetIssuesResponse(pageSize int, startCursor string) (GitlabIssuesResponse, error) {
	data, err := client.graphQLRequest(client.getIssuesQuery(pageSize, startCursor))

	if err != nil {
		return GitlabIssuesResponse{}, err
	}

	var response GitlabIssuesResponse

	err = json.Unmarshal(data, &response)

	if err != nil {
		return GitlabIssuesResponse{}, err
	}

	return response, nil
}

func (client *GitlabClient) GetProjectsResponse(pageSize int, startCursor string) (GitlabProjectsResponse, error) {
	data, err := client.graphQLRequest(client.getProjectsQuery(pageSize, startCursor))

	if err != nil {
		return GitlabProjectsResponse{}, err
	}

	var response GitlabProjectsResponse

	err = json.Unmarshal(data, &response)

	if err != nil {
		return GitlabProjectsResponse{}, err
	}

	return response, nil
}

func (client *GitlabClient) getGraphQlEndpoint() string {
	u, _ := url.JoinPath(client.apiUrl, "api/graphql")
	return u
}

func (client *GitlabClient) getUsersQuery(pageSize int, startCursor string) string {
	pagination := ""

	if pageSize > 0 {
		pagination = pagination + fmt.Sprintf(" first:%d", pageSize)
	}

	if startCursor != "" {
		pagination = pagination + fmt.Sprintf(" after:\"%s\"", startCursor)
	}

	return fmt.Sprintf(`
		query {
			users(%s) {
				nodes {
					id
					name
					username
					avatarUrl
					state
				}
				pageInfo {
					startCursor
					endCursor
					hasNextPage
					hasPreviousPage
				}
			}
		}
	`, pagination)
}

func (client *GitlabClient) getIssuesQuery(pageSize int, startCursor string) string {
	pagination := ""

	if pageSize > 0 {
		pagination = pagination + fmt.Sprintf(" first:%d", pageSize)
	}

	if startCursor != "" {
		pagination = pagination + fmt.Sprintf(" after:\"%s\"", startCursor)
	}

	return fmt.Sprintf(`
		query {
			issues(state:opened %s) {
				nodes {
					id
					iid
					assignees {
						nodes {
							id
						}
					}
					title
					type
					webUrl
					labels {
						nodes {
							id
							title
							color
							textColor
						}
					}
					projectId
					milestone {
						id
						iid
						title
						webPath
					}
				}
				pageInfo {
					startCursor
					endCursor
					hasNextPage
					hasPreviousPage
				}
			}
		}
	`, pagination)
}

func (client *GitlabClient) getProjectsQuery(pageSize int, startCursor string) string {
	pagination := ""

	if pageSize > 0 {
		pagination = pagination + fmt.Sprintf(" first:%d", pageSize)
	}

	if startCursor != "" {
		pagination = pagination + fmt.Sprintf(" after:\"%s\"", startCursor)
	}

	return fmt.Sprintf(`
		query {
			projects(%s) {
				nodes {
					id
					name
					projectMembers {
						nodes {
							user {
								id
							}
						}
					}
					milestones(sort:CREATED_DESC) {
                        nodes {
                            id
                            iid
                            title
							webPath
                        }
                    }
				}
				pageInfo {
					startCursor
					endCursor
					hasNextPage
					hasPreviousPage
				}
			}
		}
	`, pagination)
}

func (client *GitlabClient) graphQLRequest(query string) ([]byte, error) {
	data := map[string]string{
		"query": query,
	}
	jsonData, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	request, _ := http.NewRequest(http.MethodPost, client.getGraphQlEndpoint(), bytes.NewBuffer(jsonData))
	request.Header.Set("PRIVATE-TOKEN", client.token)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	httpClient := http.Client{Timeout: 30 * time.Second}

	res, err := httpClient.Do(request)

	if err != nil {
		return make([]byte, 0), err
	}

	if res.StatusCode != 200 {
		return make([]byte, 0), &InvalidResponseError{res.Status, res.StatusCode}
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return make([]byte, 0), err
	}

	return body, nil
}

func (client *GitlabClient) toDomainUser(user *GitlabUser) (*domain.User, error) {
	userId, err := client.cleanUserId(user.Id)

	if err != nil {
		return nil, err
	}

	return &domain.User{
		Id:        domain.UserId(userId),
		Name:      user.Name,
		Username:  user.Username,
		AvatarUrl: client.cleanUserAvatar(user.AvatarUrl),
		IsVisible: true,
		IsActive:  user.State == "active",
	}, nil
}

func (client *GitlabClient) cleanUserId(gid string) (uint, error) {
	id, err := strconv.ParseUint(strings.ReplaceAll(gid, "gid://gitlab/User/", ""), 10, 32)

	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func (client *GitlabClient) toDomainProject(project *GitlabProject) (*domain.Project, error) {
	projectId, err := client.cleanProjectId(project.Id)
	if err != nil {
		return nil, err
	}

	userIds := utils.Map(project.GetUserIds(), func(strId string) domain.UserId {
		val, err := client.cleanUserId(strId)
		if err != nil {
			return 0
		}
		return domain.UserId(val)
	})

	result := &domain.Project{
		Id:        domain.ProjectId(projectId),
		Name:      project.Name,
		Users:     userIds,
		IsVisible: true,
	}

	if err := result.Validate(); err != nil {
		return nil, err
	}

	return result, nil
}

func (client *GitlabClient) cleanProjectId(gid string) (uint, error) {
	id, err := strconv.ParseUint(strings.ReplaceAll(gid, "gid://gitlab/Project/", ""), 10, 32)

	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func (client *GitlabClient) toDomainIssue(
	issue *GitlabIssue,
	users []domain.User,
	projects []domain.Project,
	labels []domain.Label,
	releases []domain.Release,
	settings *domain.Settings,
) (*domain.Issue, error) {
	assignees := make([]domain.User, 0)
	for _, a := range issue.Assignees.Nodes {
		for _, user := range users {
			assigneeId, err := client.cleanUserId(a.UserId)
			if err != nil {
				continue
			}

			if user.Id == domain.UserId(assigneeId) {
				assignees = append(assignees, user)
			}
		}
	}

	projectId := domain.ProjectId(uint(issue.ProjectId))
	project := utils.Find(projects, func(p domain.Project) bool {
		return p.Id == projectId
	})
	if project == nil {
		return nil, fmt.Errorf("project not found")
	}

	issueLabels := make([]domain.Label, 0)
	for _, l := range issue.Labels.Nodes {
		for _, label := range labels {
			if domain.LabelId(l.Id) == label.Id {
				issueLabels = append(issueLabels, label)
			}
		}
	}

	release := utils.Find(releases, func(r domain.Release) bool {
		releaseId, err := client.cleanReleaseId(issue.Milestone.Id)
		return r.Id == domain.ReleaseId(releaseId) && err == nil
	})

	domainIssue := &domain.Issue{
		Id:          domain.IssueId(0),
		ExternalId:  domain.IssueExternalId(issue.Id),
		Iid:         domain.IssueIid(issue.Iid),
		Title:       issue.Title,
		IssueType:   domain.IssueType(issue.IssueType),
		Assignees:   assignees,
		WebUrl:      issue.WebUrl,
		Labels:      issueLabels,
		Project:     *project,
		Release:     release,
		TaskType:    client.getIssueTaskType(issueLabels, settings),
		EstimateDev: nil,
		EstimateQA:  nil,
	}

	if err := domainIssue.Validate(); err != nil {
		return nil, err
	}

	return domainIssue, nil
}

func (client *GitlabClient) getIssueTaskType(labels []domain.Label, settings *domain.Settings) *domain.Label {
	return utils.Find(labels, func(label domain.Label) bool {
		return utils.IndexOf(settings.TaskTypeLabels, func(id string) bool {
			return id == label.Name
		}) > -1
	})
}

func (client *GitlabClient) toDomainLabel(label *GitlabLabel) *domain.Label {
	return &domain.Label{
		Id:        domain.LabelId(label.Id),
		Name:      label.Title,
		Color:     domain.Color(label.Color),
		TextColor: domain.Color(label.TextColor),
		AltName:   nil,
	}
}

func (client *GitlabClient) toDomainRelease(milestone *GitlabMilestone, project *domain.Project) (*domain.Release, error) {
	releaseId, err := client.cleanReleaseId(milestone.Id)
	if err != nil {
		return nil, err
	}

	result := &domain.Release{
		Id:      domain.ReleaseId(releaseId),
		Iid:     domain.ReleaseIid(milestone.Iid),
		Title:   milestone.Title,
		Project: *project,
		WebPath: client.cleanReleaseWebUrl(milestone.WebPath),
	}

	if err := result.Validate(); err != nil {
		return nil, err
	}

	return result, nil
}

func (client *GitlabClient) cleanReleaseWebUrl(webUrl string) string {
	if !strings.HasPrefix(webUrl, "https://") {
		result, err := url.JoinPath(client.config.GitlabUrl, webUrl)
		if err == nil {
			return result
		} else {
			return webUrl
		}
	}
	return webUrl
}

func (client *GitlabClient) cleanUserAvatar(avatarUrl string) string {
	if !strings.HasPrefix(avatarUrl, "https://") {
		result, err := url.JoinPath(client.config.GitlabUrl, avatarUrl)
		if err == nil {
			return result
		} else {
			return avatarUrl
		}
	}
	return avatarUrl
}

func (client *GitlabClient) cleanReleaseId(gid string) (uint, error) {
	id, err := strconv.ParseUint(strings.ReplaceAll(gid, "gid://gitlab/Milestone/", ""), 10, 32)

	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func (client *GitlabClient) ExchangeCodeForToken(code string) (*interfaces.GitLabOAuthToken, error) {
	data := url.Values{}
	data.Set("client_id", client.config.GitLabAuthClientID)
	data.Set("client_secret", client.config.GitLabAuthClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", client.config.GitLabAuthRedirectURL)

	tokenURL := fmt.Sprintf("%s/oauth/token", client.config.GitlabUrl)

	req, err := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	httpClient := http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to exchange code for token: status %d, body: %s", resp.StatusCode, string(body))
	}

	var token interfaces.GitLabOAuthToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &token, nil
}

func (client *GitlabClient) RefreshToken(refreshToken string) (*interfaces.GitLabOAuthToken, error) {
	data := url.Values{}
	data.Set("client_id", client.config.GitLabAuthClientID)
	data.Set("client_secret", client.config.GitLabAuthClientSecret)
	data.Set("refresh_token", refreshToken)
	data.Set("grant_type", "refresh_token")
	data.Set("redirect_uri", client.config.GitLabAuthRedirectURL)

	tokenURL := fmt.Sprintf("%s/oauth/token", client.config.GitlabUrl)

	req, err := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	httpClient := http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to refresh token: status %d, body: %s", resp.StatusCode, string(body))
	}

	var token interfaces.GitLabOAuthToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &token, nil
}

func (client *GitlabClient) GetUserInfo(accessToken string) (*interfaces.GitLabUserInfo, error) {
	userInfoURL := fmt.Sprintf("%s/api/v4/user", client.config.GitlabUrl)

	req, err := http.NewRequest(http.MethodGet, userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/json")

	httpClient := http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: status %d, body: %s", resp.StatusCode, string(body))
	}

	var userInfo interfaces.GitLabUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}

func (client *GitlabClient) doUserRequest(userToken, method, endpoint string, body any) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userToken))
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	httpClient := http.Client{Timeout: 30 * time.Second}
	return httpClient.Do(req)
}

func (client *GitlabClient) CreateIssue(userToken string, title string, projectId uint, assigneeId *uint) (*domain.Issue, error) {
	endpoint, err := url.JoinPath(client.apiUrl, fmt.Sprintf("api/v4/projects/%d/issues", projectId))
	if err != nil {
		return nil, fmt.Errorf("failed to build issue endpoint: %w", err)
	}

	body := map[string]any{"title": title}
	if assigneeId != nil {
		body["assignee_ids"] = []uint{*assigneeId}
	}

	resp, err := client.doUserRequest(userToken, http.MethodPost, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create issue in GitLab: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitLab returned status %d: %s", resp.StatusCode, string(body))
	}

	var restIssue GitlabRestIssue
	if err := json.NewDecoder(resp.Body).Decode(&restIssue); err != nil {
		return nil, fmt.Errorf("failed to decode GitLab issue response: %w", err)
	}

	project, err := client.projectRepo.Get(domain.ProjectId(projectId))
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	assignees := make([]domain.User, 0)
	if assigneeId != nil {
		user, err := client.userRepo.Get(domain.UserId(*assigneeId))
		if err != nil {
			return nil, fmt.Errorf("failed to get assignee: %w", err)
		}
		assignees = append(assignees, *user)
	}

	return &domain.Issue{
		ExternalId: domain.IssueExternalId(fmt.Sprintf("gid://gitlab/Issue/%d", restIssue.Id)),
		Iid:        domain.IssueIid(strconv.Itoa(restIssue.Iid)),
		Title:      restIssue.Title,
		IssueType:  domain.IssueTypeIssue,
		WebUrl:     restIssue.WebUrl,
		Project:    *project,
		Assignees:  assignees,
		Labels:     []domain.Label{},
	}, nil
}

func (client *GitlabClient) UpdateIssueMilestone(userToken string, projectId uint, issueIid string, milestoneId *uint) error {
	endpoint, err := url.JoinPath(client.apiUrl, fmt.Sprintf("api/v4/projects/%d/issues/%s", projectId, issueIid))
	if err != nil {
		return fmt.Errorf("failed to build endpoint: %w", err)
	}

	var body map[string]any
	if milestoneId != nil {
		body = map[string]any{"milestone_id": *milestoneId}
	} else {
		body = map[string]any{"milestone_id": nil}
	}

	resp, err := client.doUserRequest(userToken, http.MethodPut, endpoint, body)
	if err != nil {
		return fmt.Errorf("failed to update issue milestone in GitLab: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GitLab returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (client *GitlabClient) RemoveIssueLink(userToken string, projectId uint, issueIid string, targetProjectId uint, targetIssueIid string) error {
	listEndpoint, err := url.JoinPath(client.apiUrl, fmt.Sprintf("api/v4/projects/%d/issues/%s/links", projectId, issueIid))
	if err != nil {
		return fmt.Errorf("failed to build endpoint: %w", err)
	}

	resp, err := client.doUserRequest(userToken, http.MethodGet, listEndpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to get issue links from GitLab: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GitLab returned status %d: %s", resp.StatusCode, string(body))
	}

	var links []GitlabIssueLinkItem
	if err := json.NewDecoder(resp.Body).Decode(&links); err != nil {
		return fmt.Errorf("failed to decode issue links response: %w", err)
	}

	var linkId uint
	for _, link := range links {
		if link.ProjectId == targetProjectId && fmt.Sprintf("%d", link.Iid) == targetIssueIid {
			linkId = link.IssueLinkId
			break
		}
	}
	if linkId == 0 {
		return nil
	}

	deleteEndpoint, err := url.JoinPath(client.apiUrl, fmt.Sprintf("api/v4/projects/%d/issues/%s/links/%d", projectId, issueIid, linkId))
	if err != nil {
		return fmt.Errorf("failed to build delete endpoint: %w", err)
	}

	delResp, err := client.doUserRequest(userToken, http.MethodDelete, deleteEndpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to delete issue link in GitLab: %w", err)
	}
	defer delResp.Body.Close()

	if delResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(delResp.Body)
		return fmt.Errorf("GitLab returned status %d: %s", delResp.StatusCode, string(body))
	}

	return nil
}

func (client *GitlabClient) AddIssueLink(userToken string, projectId uint, issueIid string, targetProjectId uint, targetIssueIid string) error {
	endpoint, err := url.JoinPath(client.apiUrl, fmt.Sprintf("api/v4/projects/%d/issues/%s/links", projectId, issueIid))
	if err != nil {
		return fmt.Errorf("failed to build endpoint: %w", err)
	}

	body := map[string]any{
		"target_project_id": targetProjectId,
		"target_issue_iid":  targetIssueIid,
	}

	resp, err := client.doUserRequest(userToken, http.MethodPost, endpoint, body)
	if err != nil {
		return fmt.Errorf("failed to add issue link in GitLab: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GitLab returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (client *GitlabClient) AddIssueLabel(userToken string, projectId uint, issueIid string, labelName string) error {
	endpoint, err := url.JoinPath(client.apiUrl, fmt.Sprintf("api/v4/projects/%d/issues/%s", projectId, issueIid))
	if err != nil {
		return fmt.Errorf("failed to build endpoint: %w", err)
	}

	resp, err := client.doUserRequest(userToken, http.MethodPut, endpoint, map[string]any{"add_labels": labelName})
	if err != nil {
		return fmt.Errorf("failed to add issue label in GitLab: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GitLab returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (client *GitlabClient) doSystemRequest(method, endpoint string, body any) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("PRIVATE-TOKEN", client.token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	httpClient := http.Client{Timeout: 30 * time.Second}
	return httpClient.Do(req)
}

func (client *GitlabClient) GetLinkedIssues(projectId uint, issueIid string) ([]domain.IssueExternalId, error) {
	endpoint, err := url.JoinPath(client.apiUrl, fmt.Sprintf("api/v4/projects/%d/issues/%s/links", projectId, issueIid))
	if err != nil {
		return nil, fmt.Errorf("failed to build endpoint: %w", err)
	}

	resp, err := client.doSystemRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get linked issues from GitLab: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitLab returned status %d: %s", resp.StatusCode, string(body))
	}

	var links []GitlabIssueLinkItem
	if err := json.NewDecoder(resp.Body).Decode(&links); err != nil {
		return nil, fmt.Errorf("failed to decode linked issues response: %w", err)
	}

	result := make([]domain.IssueExternalId, 0, len(links))
	for _, link := range links {
		result = append(result, domain.IssueExternalId(fmt.Sprintf("gid://gitlab/Issue/%d", link.Id)))
	}
	return result, nil
}

func (client *GitlabClient) RemoveIssueLabel(userToken string, projectId uint, issueIid string, labelName string) error {
	endpoint, err := url.JoinPath(client.apiUrl, fmt.Sprintf("api/v4/projects/%d/issues/%s", projectId, issueIid))
	if err != nil {
		return fmt.Errorf("failed to build endpoint: %w", err)
	}

	resp, err := client.doUserRequest(userToken, http.MethodPut, endpoint, map[string]any{"remove_labels": labelName})
	if err != nil {
		return fmt.Errorf("failed to remove issue label in GitLab: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GitLab returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
