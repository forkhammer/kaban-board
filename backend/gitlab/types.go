package gitlab

type GitlabUsersResponse struct {
	Data struct {
		Users struct {
			Nodes    []GitlabUser   `json:"nodes"`
			PageInfo GitlabPageInfo `json:"pageInfo"`
		} `json:"users"`
	} `json:"data"`
}

type GitlabUser struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatarUrl"`
}

type GitlabAssignee struct {
	UserId string `json:"id"`
}

type GitlabIssue struct {
	Id        string `json:"id"`
	Iid       string `json:"iid"`
	Title     string `json:"title"`
	IssueType string `json:"type"`
	Assignees struct {
		Nodes []GitlabAssignee `json:"nodes"`
	} `json:"assignees"`
	WebUrl string `json:"webUrl"`
	Labels struct {
		Nodes []GitlabLabel `json:"nodes"`
	} `json:"labels"`
	ProjectId   int             `json:"projectId"`
	ProjectName *string         `json:"projectName"`
	Milestone   GitlabMilestone `json:"milestone"`
}

type GitlabLabel struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Color     string `json:"color"`
	TextColor string `json:"textColor"`
}

type GitlabIssuesResponse struct {
	Data struct {
		Issues struct {
			Nodes    []GitlabIssue  `json:"nodes"`
			PageInfo GitlabPageInfo `json:"pageInfo"`
		} `json:"issues"`
	} `json:"data"`
}

type GitlabPageInfo struct {
	StartCursor     string `json:"startCursor"`
	EndCursor       string `json:"endCursor"`
	HasNextPage     bool   `json:"hasNextPage"`
	HasPreviousPage bool   `json:"hasPreviousPage"`
}

type GitlabProject struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	ProjectMembers struct {
		Nodes []GitlabProjectMember `json:"nodes"`
	} `json:"projectMembers"`
}

func (p *GitlabProject) GetUserIds() []string {
	result := make([]string, 0)

	for i := range p.ProjectMembers.Nodes {
		pm := &p.ProjectMembers.Nodes[i]
		result = append(result, pm.User.Id)
	}

	return result
}

type GitlabProjectsResponse struct {
	Data struct {
		Projects struct {
			Nodes    []GitlabProject `json:"nodes"`
			PageInfo GitlabPageInfo  `json:"pageInfo"`
		} `json:"projects"`
	} `json:"data"`
}

type GitlabProjectMember struct {
	User struct {
		Id string `json:"id"`
	} `json:"user"`
}

type GitlabMilestone struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	WebPath string `json:"webPath"`
}
