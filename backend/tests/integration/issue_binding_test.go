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

type bindingPageResult struct {
	Count      int              `json:"count"`
	Pages      int              `json:"pages"`
	Page       int              `json:"page"`
	StartIndex int              `json:"start_index"`
	EndIndex   int              `json:"end_index"`
	Results    []map[string]any `json:"results"`
}

func TestGetBindings(t *testing.T) {
	testutil.Run(t, "returns empty list when no bindings", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/binding", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result bindingPageResult
		testutil.ParseBody(t, resp, &result)

		assert.Empty(t, result.Results)
		assert.Equal(t, 0, result.Count)
		assert.Equal(t, 0, result.Pages)
	})

	testutil.Run(t, "returns bindings with issue details", func(t *testing.T) {
		team, err := factories.NewTeamFactory().Create(nil)
		require.NoError(t, err)

		project, err := factories.NewProjectFactory().Create(map[string]any{"team": team})
		require.NoError(t, err)

		issue, err := factories.NewIssueFactory().Create(map[string]any{"project": *project})
		require.NoError(t, err)

		sprint, err := factories.NewSprintFactory().Create(map[string]any{"team": *team})
		require.NoError(t, err)

		binding, err := factories.NewIssueBindingFactory().Create(map[string]any{
			"sprint": sprint,
			"issue":  issue,
		})
		require.NoError(t, err)

		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/binding", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result bindingPageResult
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result.Results, 1)
		assert.Equal(t, float64(binding.Id), result.Results[0]["bindingId"])
		assert.Equal(t, issue.Title, result.Results[0]["title"])
		assert.Equal(t, float64(project.Id), result.Results[0]["projectId"])
		assert.Equal(t, "backlog", result.Results[0]["bindStatus"])
		assert.NotNil(t, result.Results[0]["sprint"])
	})

	testutil.Run(t, "filters by sprint", func(t *testing.T) {
		team, err := factories.NewTeamFactory().Create(nil)
		require.NoError(t, err)

		project, err := factories.NewProjectFactory().Create(map[string]any{"team": team})
		require.NoError(t, err)

		sprint1, err := factories.NewSprintFactory().Create(map[string]any{
			"team":  *team,
			"title": "Sprint One",
		})
		require.NoError(t, err)

		sprint2, err := factories.NewSprintFactory().Create(map[string]any{
			"team":  *team,
			"title": "Sprint Two",
		})
		require.NoError(t, err)

		issue1, err := factories.NewIssueFactory().Create(map[string]any{
			"project": *project,
			"title":   "Issue for Sprint 1",
		})
		require.NoError(t, err)

		issue2, err := factories.NewIssueFactory().Create(map[string]any{
			"project": *project,
			"title":   "Issue for Sprint 2",
		})
		require.NoError(t, err)

		_, err = factories.NewIssueBindingFactory().Create(map[string]any{
			"sprint": sprint1,
			"issue":  issue1,
		})
		require.NoError(t, err)

		_, err = factories.NewIssueBindingFactory().Create(map[string]any{
			"sprint": sprint2,
			"issue":  issue2,
		})
		require.NoError(t, err)

		path := fmt.Sprintf("/api/binding?sprint=%d", sprint1.Id)
		resp := testutil.ReqJSON(t, suite.Server, "GET", path, nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result bindingPageResult
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result.Results, 1)
		assert.Equal(t, "Issue for Sprint 1", result.Results[0]["title"])
	})

	testutil.Run(t, "filters by issue", func(t *testing.T) {
		team, err := factories.NewTeamFactory().Create(nil)
		require.NoError(t, err)

		project, err := factories.NewProjectFactory().Create(map[string]any{"team": team})
		require.NoError(t, err)

		sprint, err := factories.NewSprintFactory().Create(map[string]any{"team": *team})
		require.NoError(t, err)

		issue1, err := factories.NewIssueFactory().Create(map[string]any{
			"project": *project,
			"title":   "Target Issue",
		})
		require.NoError(t, err)

		issue2, err := factories.NewIssueFactory().Create(map[string]any{
			"project": *project,
			"title":   "Other Issue",
		})
		require.NoError(t, err)

		_, err = factories.NewIssueBindingFactory().Create(map[string]any{
			"sprint": sprint,
			"issue":  issue1,
		})
		require.NoError(t, err)

		_, err = factories.NewIssueBindingFactory().Create(map[string]any{
			"sprint": sprint,
			"issue":  issue2,
		})
		require.NoError(t, err)

		path := fmt.Sprintf("/api/binding?issue=%d", issue1.Id)
		resp := testutil.ReqJSON(t, suite.Server, "GET", path, nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result bindingPageResult
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result.Results, 1)
		assert.Equal(t, "Target Issue", result.Results[0]["title"])
	})

	testutil.Run(t, "supports pagination", func(t *testing.T) {
		team, err := factories.NewTeamFactory().Create(nil)
		require.NoError(t, err)

		project, err := factories.NewProjectFactory().Create(map[string]any{"team": team})
		require.NoError(t, err)

		sprint, err := factories.NewSprintFactory().Create(map[string]any{"team": *team})
		require.NoError(t, err)

		for i := range 3 {
			issue, err := factories.NewIssueFactory().Create(map[string]any{
				"project": *project,
				"title":   fmt.Sprintf("Binding %d", i+1),
			})
			require.NoError(t, err)

			_, err = factories.NewIssueBindingFactory().Create(map[string]any{
				"sprint": sprint,
				"issue":  issue,
			})
			require.NoError(t, err)
		}

		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/binding?page=1&limit=2", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result bindingPageResult
		testutil.ParseBody(t, resp, &result)

		assert.Equal(t, 3, result.Count)
		assert.Equal(t, 2, result.Pages)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 0, result.StartIndex)
		assert.Equal(t, 2, result.EndIndex)
		require.Len(t, result.Results, 2)
	})
}

func TestGetBinding(t *testing.T) {
	testutil.Run(t, "returns existing binding", func(t *testing.T) {
		team, err := factories.NewTeamFactory().Create(nil)
		require.NoError(t, err)

		project, err := factories.NewProjectFactory().Create(map[string]any{"team": team})
		require.NoError(t, err)

		issue, err := factories.NewIssueFactory().Create(map[string]any{"project": *project})
		require.NoError(t, err)

		sprint, err := factories.NewSprintFactory().Create(map[string]any{"team": *team})
		require.NoError(t, err)

		binding, err := factories.NewIssueBindingFactory().Create(map[string]any{
			"sprint": sprint,
			"issue":  issue,
		})
		require.NoError(t, err)

		path := fmt.Sprintf("/api/binding/%d", binding.Id)
		resp := testutil.ReqJSON(t, suite.Server, "GET", path, nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Equal(t, float64(binding.Id), result["bindingId"])
		assert.NotNil(t, result["can_update"])
		assert.NotNil(t, result["can_manage"])
		assert.NotNil(t, result["sprint"])
	})

	testutil.Run(t, "returns 404 for nonexistent binding", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/binding/999999", nil, nil)

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestCreateBinding(t *testing.T) {
	testutil.Run(t, "creates issue and binding when authenticated", func(t *testing.T) {
		team, err := factories.NewTeamFactory().Create(nil)
		require.NoError(t, err)

		project, err := factories.NewProjectFactory().Create(map[string]any{"team": team})
		require.NoError(t, err)

		sprint, err := factories.NewSprintFactory().Create(map[string]any{"team": *team})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		body := map[string]any{
			"title":           "Test Issue",
			"project":         project.Id,
			"sprint":          sprint.Id,
			"createInTracker": false,
		}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/binding", body, &token)

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.NotZero(t, result["bindingId"])
		assert.Equal(t, "Test Issue", result["title"])
		assert.Equal(t, "backlog", result["bindStatus"])
		assert.NotNil(t, result["sprint"])
	})

	testutil.Run(t, "returns 401 without auth", func(t *testing.T) {
		body := map[string]any{
			"title":           "Unauthorized",
			"project":         uint(1),
			"sprint":          uint(1),
			"createInTracker": false,
		}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/binding", body, nil)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	testutil.Run(t, "returns 400 for missing required fields", func(t *testing.T) {
		token := testutil.GetAdminToken(t, suite.Server)

		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/binding", map[string]any{}, &token)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestSaveBinding(t *testing.T) {
	testutil.Run(t, "saves binding when authenticated", func(t *testing.T) {
		team, err := factories.NewTeamFactory().Create(nil)
		require.NoError(t, err)

		project, err := factories.NewProjectFactory().Create(map[string]any{"team": team})
		require.NoError(t, err)

		issue, err := factories.NewIssueFactory().Create(map[string]any{"project": *project})
		require.NoError(t, err)

		sprint, err := factories.NewSprintFactory().Create(map[string]any{"team": *team})
		require.NoError(t, err)

		binding, err := factories.NewIssueBindingFactory().Create(map[string]any{
			"sprint": sprint,
			"issue":  issue,
		})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		// GET binding to read current version
		getPath := fmt.Sprintf("/api/binding/%d", binding.Id)
		getResp := testutil.ReqJSON(t, suite.Server, "GET", getPath, nil, nil)
		var current map[string]any
		testutil.ParseBody(t, getResp, &current)

		body := map[string]any{
			"bindStatus": "in_progress",
			"version":    current["version"],
		}
		path := fmt.Sprintf("/api/binding/%d", binding.Id)
		resp := testutil.ReqJSON(t, suite.Server, "PUT", path, body, &token)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Equal(t, "in_progress", result["bindStatus"])
	})

	testutil.Run(t, "returns 401 without auth", func(t *testing.T) {
		body := map[string]any{
			"title": "Test",
		}
		resp := testutil.ReqJSON(t, suite.Server, "PUT", "/api/binding/1", body, nil)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	testutil.Run(t, "returns 409 on version conflict", func(t *testing.T) {
		team, err := factories.NewTeamFactory().Create(nil)
		require.NoError(t, err)

		project, err := factories.NewProjectFactory().Create(map[string]any{"team": team})
		require.NoError(t, err)

		issue, err := factories.NewIssueFactory().Create(map[string]any{"project": *project})
		require.NoError(t, err)

		sprint, err := factories.NewSprintFactory().Create(map[string]any{"team": *team})
		require.NoError(t, err)

		binding, err := factories.NewIssueBindingFactory().Create(map[string]any{
			"sprint": sprint,
			"issue":  issue,
		})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		body := map[string]any{
			"title":   "Conflict",
			"version": uint(0),
		}
		path := fmt.Sprintf("/api/binding/%d", binding.Id)
		resp := testutil.ReqJSON(t, suite.Server, "PUT", path, body, &token)

		require.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	testutil.Run(t, "returns 404 for nonexistent binding", func(t *testing.T) {
		token := testutil.GetAdminToken(t, suite.Server)

		body := map[string]any{
			"title":   "Not Found",
			"version": uint(1),
		}
		resp := testutil.ReqJSON(t, suite.Server, "PUT", "/api/binding/999999", body, &token)

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestDeleteBinding(t *testing.T) {
	testutil.Run(t, "deletes binding as admin", func(t *testing.T) {
		team, err := factories.NewTeamFactory().Create(nil)
		require.NoError(t, err)

		project, err := factories.NewProjectFactory().Create(map[string]any{"team": team})
		require.NoError(t, err)

		issue, err := factories.NewIssueFactory().Create(map[string]any{"project": *project})
		require.NoError(t, err)

		sprint, err := factories.NewSprintFactory().Create(map[string]any{"team": *team})
		require.NoError(t, err)

		binding, err := factories.NewIssueBindingFactory().Create(map[string]any{
			"sprint": sprint,
			"issue":  issue,
		})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		path := fmt.Sprintf("/api/binding/%d", binding.Id)
		resp := testutil.ReqJSON(t, suite.Server, "DELETE", path, nil, &token)

		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	testutil.Run(t, "returns 403 as viewer", func(t *testing.T) {
		team, err := factories.NewTeamFactory().Create(nil)
		require.NoError(t, err)

		project, err := factories.NewProjectFactory().Create(map[string]any{"team": team})
		require.NoError(t, err)

		issue, err := factories.NewIssueFactory().Create(map[string]any{"project": *project})
		require.NoError(t, err)

		sprint, err := factories.NewSprintFactory().Create(map[string]any{"team": *team})
		require.NoError(t, err)

		binding, err := factories.NewIssueBindingFactory().Create(map[string]any{
			"sprint": sprint,
			"issue":  issue,
		})
		require.NoError(t, err)

		token := testutil.GetViewerToken(t, suite.Server)

		path := fmt.Sprintf("/api/binding/%d", binding.Id)
		resp := testutil.ReqJSON(t, suite.Server, "DELETE", path, nil, &token)

		require.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	testutil.Run(t, "returns 401 without auth", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "DELETE", "/api/binding/1", nil, nil)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestCopyBinding(t *testing.T) {
	testutil.Run(t, "copies binding to another sprint as admin", func(t *testing.T) {
		team, err := factories.NewTeamFactory().Create(nil)
		require.NoError(t, err)

		project, err := factories.NewProjectFactory().Create(map[string]any{"team": team})
		require.NoError(t, err)

		issue, err := factories.NewIssueFactory().Create(map[string]any{"project": *project})
		require.NoError(t, err)

		sprintA, err := factories.NewSprintFactory().Create(map[string]any{
			"team":  *team,
			"title": "Sprint A",
		})
		require.NoError(t, err)

		sprintB, err := factories.NewSprintFactory().Create(map[string]any{
			"team":  *team,
			"title": "Sprint B",
		})
		require.NoError(t, err)

		binding, err := factories.NewIssueBindingFactory().Create(map[string]any{
			"sprint": sprintA,
			"issue":  issue,
		})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		body := map[string]any{
			"sprint_id": sprintB.Id,
		}
		path := fmt.Sprintf("/api/binding/%d/copy", binding.Id)
		resp := testutil.ReqJSON(t, suite.Server, "POST", path, body, &token)

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.GreaterOrEqual(t, result["bindingId"], float64(binding.Id+1))
		assert.NotEqual(t, float64(binding.Id), result["bindingId"])

		sprintObj, ok := result["sprint"].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, float64(sprintB.Id), sprintObj["id"])
	})

	testutil.Run(t, "returns 403 as viewer", func(t *testing.T) {
		team, err := factories.NewTeamFactory().Create(nil)
		require.NoError(t, err)

		project, err := factories.NewProjectFactory().Create(map[string]any{"team": team})
		require.NoError(t, err)

		issue, err := factories.NewIssueFactory().Create(map[string]any{"project": *project})
		require.NoError(t, err)

		sprint, err := factories.NewSprintFactory().Create(map[string]any{"team": *team})
		require.NoError(t, err)

		binding, err := factories.NewIssueBindingFactory().Create(map[string]any{
			"sprint": sprint,
			"issue":  issue,
		})
		require.NoError(t, err)

		token := testutil.GetViewerToken(t, suite.Server)

		body := map[string]any{
			"sprint_id": sprint.Id,
		}
		path := fmt.Sprintf("/api/binding/%d/copy", binding.Id)
		resp := testutil.ReqJSON(t, suite.Server, "POST", path, body, &token)

		require.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	testutil.Run(t, "returns 401 without auth", func(t *testing.T) {
		body := map[string]any{
			"sprint_id": uint(1),
		}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/binding/1/copy", body, nil)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestMoveBinding(t *testing.T) {
	testutil.Run(t, "moves binding to another sprint as admin", func(t *testing.T) {
		team, err := factories.NewTeamFactory().Create(nil)
		require.NoError(t, err)

		project, err := factories.NewProjectFactory().Create(map[string]any{"team": team})
		require.NoError(t, err)

		issue, err := factories.NewIssueFactory().Create(map[string]any{"project": *project})
		require.NoError(t, err)

		sprintA, err := factories.NewSprintFactory().Create(map[string]any{
			"team":  *team,
			"title": "Sprint A",
		})
		require.NoError(t, err)

		sprintB, err := factories.NewSprintFactory().Create(map[string]any{
			"team":  *team,
			"title": "Sprint B",
		})
		require.NoError(t, err)

		binding, err := factories.NewIssueBindingFactory().Create(map[string]any{
			"sprint": sprintA,
			"issue":  issue,
		})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		body := map[string]any{
			"sprint_id": sprintB.Id,
		}
		path := fmt.Sprintf("/api/binding/%d/move", binding.Id)
		resp := testutil.ReqJSON(t, suite.Server, "POST", path, body, &token)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Equal(t, float64(binding.Id), result["bindingId"])

		sprintObj, ok := result["sprint"].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, float64(sprintB.Id), sprintObj["id"])
	})

	testutil.Run(t, "returns 403 as viewer", func(t *testing.T) {
		team, err := factories.NewTeamFactory().Create(nil)
		require.NoError(t, err)

		project, err := factories.NewProjectFactory().Create(map[string]any{"team": team})
		require.NoError(t, err)

		issue, err := factories.NewIssueFactory().Create(map[string]any{"project": *project})
		require.NoError(t, err)

		sprint, err := factories.NewSprintFactory().Create(map[string]any{"team": *team})
		require.NoError(t, err)

		binding, err := factories.NewIssueBindingFactory().Create(map[string]any{
			"sprint": sprint,
			"issue":  issue,
		})
		require.NoError(t, err)

		token := testutil.GetViewerToken(t, suite.Server)

		body := map[string]any{
			"sprint_id": sprint.Id,
		}
		path := fmt.Sprintf("/api/binding/%d/move", binding.Id)
		resp := testutil.ReqJSON(t, suite.Server, "POST", path, body, &token)

		require.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	testutil.Run(t, "returns 401 without auth", func(t *testing.T) {
		body := map[string]any{
			"sprint_id": uint(1),
		}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/binding/1/move", body, nil)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestSaveOrdering(t *testing.T) {
	testutil.Run(t, "saves ordering as admin", func(t *testing.T) {
		team, err := factories.NewTeamFactory().Create(nil)
		require.NoError(t, err)

		project, err := factories.NewProjectFactory().Create(map[string]any{"team": team})
		require.NoError(t, err)

		sprint, err := factories.NewSprintFactory().Create(map[string]any{"team": *team})
		require.NoError(t, err)

		issue1, err := factories.NewIssueFactory().Create(map[string]any{"project": *project})
		require.NoError(t, err)

		issue2, err := factories.NewIssueFactory().Create(map[string]any{"project": *project})
		require.NoError(t, err)

		binding1, err := factories.NewIssueBindingFactory().Create(map[string]any{
			"sprint": sprint,
			"issue":  issue1,
		})
		require.NoError(t, err)

		binding2, err := factories.NewIssueBindingFactory().Create(map[string]any{
			"sprint": sprint,
			"issue":  issue2,
		})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		body := []map[string]any{
			{"id": binding1.Id, "order": "1"},
			{"id": binding2.Id, "order": "2"},
		}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/binding/save_ordering", body, &token)

		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	testutil.Run(t, "returns 403 as viewer", func(t *testing.T) {
		token := testutil.GetViewerToken(t, suite.Server)

		body := []map[string]any{
			{"id": uint(1), "order": "1"},
		}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/binding/save_ordering", body, &token)

		require.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	testutil.Run(t, "returns 401 without auth", func(t *testing.T) {
		body := []map[string]any{
			{"id": uint(1), "order": "1"},
		}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/binding/save_ordering", body, nil)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
