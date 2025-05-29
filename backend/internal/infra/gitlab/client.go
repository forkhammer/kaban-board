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

	"main/internal/app/interfaces"
	domain "main/internal/domain/models"
)

type GitlabClient struct {
	interfaces.TaskTracker
	apiUrl string
	token  string
}

func NewGitlabClient(apiUrl string, token string) *GitlabClient {
	return &GitlabClient{
		apiUrl: apiUrl,
		token:  token,
	}
}

func (client *GitlabClient) GetProjects() ([]domain.Project, error) {
	return nil, nil
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
	return nil, nil
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
		AvatarUrl: user.AvatarUrl,
		IsVisible: true,
	}, nil
}

func (client *GitlabClient) cleanUserId(gid string) (uint, error) {
	id, err := strconv.ParseUint(strings.ReplaceAll(gid, "gid://gitlab/User/", ""), 10, 32)

	if err != nil {
		return 0, err
	}

	return uint(id), nil
}
