package integration

import (
	"fmt"
	"main/internal/testutil"
	"main/internal/testutil/factories"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type issuePageResult struct {
	Count      int              `json:"count"`
	Pages      int              `json:"pages"`
	Page       int              `json:"page"`
	StartIndex int              `json:"start_index"`
	EndIndex   int              `json:"end_index"`
	Results    []map[string]any `json:"results"`
}

type bindIssueResponse struct {
	Results []map[string]any `json:"results"`
	Errors  []string         `json:"errors"`
}

func TestGetIssues(t *testing.T) {
	testutil.Run(t, "returns empty list when no issues", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/issue", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result issuePageResult
		testutil.ParseBody(t, resp, &result)

		assert.Empty(t, result.Results)
		assert.Equal(t, 0, result.Count)
		assert.Equal(t, 0, result.Pages)
	})

	testutil.Run(t, "returns all issues", func(t *testing.T) {
		projectFactory := factories.NewProjectFactory()
		project, err := projectFactory.Create(nil)
		require.NoError(t, err)

		issueFactory := factories.NewIssueFactory()
		issue1, err := issueFactory.Create(map[string]any{
			"project": *project,
			"title":   "Issue One",
		})
		require.NoError(t, err)

		issue2, err := issueFactory.Create(map[string]any{
			"project": *project,
			"title":   "Issue Two",
		})
		require.NoError(t, err)

		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/issue", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result issuePageResult
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result.Results, 2)

		ids := make([]string, len(result.Results))
		for i, r := range result.Results {
			ids[i] = r["id"].(string)
		}

		assert.Contains(t, ids, fmt.Sprintf("%d", issue1.Id))
		assert.Contains(t, ids, fmt.Sprintf("%d", issue2.Id))
	})

	testutil.Run(t, "filters by project", func(t *testing.T) {
		projectFactory := factories.NewProjectFactory()
		projectA, err := projectFactory.Create(map[string]any{"name": "Project A"})
		require.NoError(t, err)

		projectB, err := projectFactory.Create(map[string]any{"name": "Project B"})
		require.NoError(t, err)

		issueFactory := factories.NewIssueFactory()
		_, err = issueFactory.Create(map[string]any{
			"project": *projectA,
			"title":   "Issue for Project A",
		})
		require.NoError(t, err)

		_, err = issueFactory.Create(map[string]any{
			"project": *projectB,
			"title":   "Issue for Project B",
		})
		require.NoError(t, err)

		path := fmt.Sprintf("/api/issue?project=%d", projectA.Id)
		resp := testutil.ReqJSON(t, suite.Server, "GET", path, nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result issuePageResult
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result.Results, 1)
		assert.Equal(t, "Issue for Project A", result.Results[0]["title"])
	})

	testutil.Run(t, "filters by search", func(t *testing.T) {
		projectFactory := factories.NewProjectFactory()
		project, err := projectFactory.Create(nil)
		require.NoError(t, err)

		issueFactory := factories.NewIssueFactory()
		_, err = issueFactory.Create(map[string]any{
			"project": *project,
			"title":   "UniqueSearchable First",
		})
		require.NoError(t, err)

		_, err = issueFactory.Create(map[string]any{
			"project": *project,
			"title":   "Other Issue",
		})
		require.NoError(t, err)

		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/issue?search=First", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result issuePageResult
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result.Results, 1)
		assert.Equal(t, "UniqueSearchable First", result.Results[0]["title"])
	})

	testutil.Run(t, "supports pagination", func(t *testing.T) {
		projectFactory := factories.NewProjectFactory()
		project, err := projectFactory.Create(nil)
		require.NoError(t, err)

		issueFactory := factories.NewIssueFactory()
		for i := range 5 {
			_, err = issueFactory.Create(map[string]any{
				"project": *project,
				"title":   fmt.Sprintf("Issue %d", i+1),
			})
			require.NoError(t, err)
		}

		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/issue?page=1&limit=2", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result issuePageResult
		testutil.ParseBody(t, resp, &result)

		assert.Equal(t, 5, result.Count)
		assert.Equal(t, 3, result.Pages)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 0, result.StartIndex)
		assert.Equal(t, 2, result.EndIndex)
		require.Len(t, result.Results, 2)
	})
}

func TestGetIssue(t *testing.T) {
	testutil.Run(t, "returns existing issue", func(t *testing.T) {
		projectFactory := factories.NewProjectFactory()
		project, err := projectFactory.Create(map[string]any{"name": "Test Project"})
		require.NoError(t, err)

		issueFactory := factories.NewIssueFactory()
		issue, err := issueFactory.Create(map[string]any{
			"project": *project,
			"title":   "My Issue",
		})
		require.NoError(t, err)

		path := fmt.Sprintf("/api/issue/%d", issue.Id)
		resp := testutil.ReqJSON(t, suite.Server, "GET", path, nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Equal(t, fmt.Sprintf("%d", issue.Id), result["id"])
		assert.Equal(t, "My Issue", result["title"])
		assert.Equal(t, float64(project.Id), result["projectId"])
		assert.Equal(t, "Test Project", result["projectName"])
	})

	testutil.Run(t, "returns 404 for nonexistent issue", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/issue/999999", nil, nil)

		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.NotEmpty(t, result["error"])
	})
}

func TestBindIssue(t *testing.T) {
	testutil.Run(t, "binds issues to sprint as admin", func(t *testing.T) {
		projectFactory := factories.NewProjectFactory()
		project, err := projectFactory.Create(nil)
		require.NoError(t, err)

		issueFactory := factories.NewIssueFactory()
		issue, err := issueFactory.Create(map[string]any{"project": *project})
		require.NoError(t, err)

		teamFactory := factories.NewTeamFactory()
		team, err := teamFactory.Create(nil)
		require.NoError(t, err)

		sprintFactory := factories.NewSprintFactory()
		sprint, err := sprintFactory.Create(map[string]any{"team": *team})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		body := map[string]any{
			"issue_ids":   []uint{uint(issue.Id)},
			"sprint_id":   sprint.Id,
			"assignee_id": 0,
		}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/issue/bind", body, &token)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result bindIssueResponse
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result.Results, 1)
	})

	testutil.Run(t, "returns multiple bindings for same issue", func(t *testing.T) {
		projectFactory := factories.NewProjectFactory()
		project, err := projectFactory.Create(nil)
		require.NoError(t, err)

		issueFactory := factories.NewIssueFactory()
		issue, err := issueFactory.Create(map[string]any{"project": *project})
		require.NoError(t, err)

		teamFactory := factories.NewTeamFactory()
		team, err := teamFactory.Create(nil)
		require.NoError(t, err)

		sprintFactory := factories.NewSprintFactory()
		sprint, err := sprintFactory.Create(map[string]any{"team": *team})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		body := map[string]any{
			"issue_ids":   []uint{uint(issue.Id)},
			"sprint_id":   sprint.Id,
			"assignee_id": 0,
		}

		// First bind — should succeed
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/issue/bind", body, &token)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		// Second bind with the same issue — also succeeds, creating a second binding
		resp = testutil.ReqJSON(t, suite.Server, "POST", "/api/issue/bind", body, &token)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result bindIssueResponse
		testutil.ParseBody(t, resp, &result)

		assert.Empty(t, result.Errors)
		require.Len(t, result.Results, 1)
	})

	testutil.Run(t, "returns 401 without auth", func(t *testing.T) {
		body := map[string]any{
			"issue_ids":   []uint{1},
			"sprint_id":   1,
			"assignee_id": 0,
		}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/issue/bind", body, nil)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	testutil.Run(t, "returns 400 for invalid body", func(t *testing.T) {
		token := testutil.GetAdminToken(t, suite.Server)

		body := map[string]any{
			"issue_ids": "not-an-array",
		}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/issue/bind", body, &token)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
