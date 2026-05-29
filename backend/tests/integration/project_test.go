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

func TestGetProjects(t *testing.T) {
	testutil.Run(t, "returns empty list when no projects", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/projects", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Empty(t, result)
	})

	testutil.Run(t, "returns all projects", func(t *testing.T) {
		factory := factories.NewProjectFactory()
		project1, err := factory.Create(map[string]any{"name": "Project Alpha"})
		require.NoError(t, err)

		project2, err := factory.Create(map[string]any{"name": "Project Beta"})
		require.NoError(t, err)

		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/projects", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result, 2)

		ids := make([]uint, len(result))
		for i, r := range result {
			ids[i] = uint(r["id"].(float64))
		}

		assert.Contains(t, ids, uint(project1.Id))
		assert.Contains(t, ids, uint(project2.Id))
	})

	testutil.Run(t, "filters by team_id", func(t *testing.T) {
		teamFactory := factories.NewTeamFactory()
		team, err := teamFactory.Create(map[string]any{"title": "Dev Team"})
		require.NoError(t, err)

		projectFactory := factories.NewProjectFactory()
		project1, err := projectFactory.Create(map[string]any{
			"name": "Project With Team",
			"team": team,
		})
		require.NoError(t, err)

		_, err = projectFactory.Create(map[string]any{
			"name": "Project Without Team",
		})
		require.NoError(t, err)

		path := fmt.Sprintf("/api/projects?team_id=%d", team.Id)
		resp := testutil.ReqJSON(t, suite.Server, "GET", path, nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result, 1)
		assert.Equal(t, float64(project1.Id), result[0]["id"])
	})

	testutil.Run(t, "filters by search", func(t *testing.T) {
		projectFactory := factories.NewProjectFactory()
		_, err := projectFactory.Create(map[string]any{"name": "Alpha Project"})
		require.NoError(t, err)

		_, err = projectFactory.Create(map[string]any{"name": "Beta Project"})
		require.NoError(t, err)

		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/projects?search=Alpha", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result, 1)
		assert.Equal(t, "Alpha Project", result[0]["name"])
	})
}

func TestGetProject(t *testing.T) {
	testutil.Run(t, "returns existing project", func(t *testing.T) {
		factory := factories.NewProjectFactory()
		project, err := factory.Create(map[string]any{"name": "Test Project"})
		require.NoError(t, err)

		path := fmt.Sprintf("/api/projects/%d", project.Id)
		resp := testutil.ReqJSON(t, suite.Server, "GET", path, nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Equal(t, float64(project.Id), result["id"])
		assert.Equal(t, "Test Project", result["name"])
	})

	testutil.Run(t, "returns 404 for nonexistent project", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/projects/999999", nil, nil)

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestSetProjectTeam(t *testing.T) {
	testutil.Run(t, "sets team as admin", func(t *testing.T) {
		projectFactory := factories.NewProjectFactory()
		project, err := projectFactory.Create(map[string]any{"name": "No Team Project"})
		require.NoError(t, err)

		teamFactory := factories.NewTeamFactory()
		team, err := teamFactory.Create(map[string]any{"title": "New Team"})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		body := map[string]any{
			"team_id": team.Id,
		}
		path := fmt.Sprintf("/api/projects/%d/set_team", project.Id)
		resp := testutil.ReqJSON(t, suite.Server, "POST", path, body, &token)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		teamIdFloat, ok := result["team_id"].(float64)
		require.True(t, ok)
		assert.Equal(t, float64(team.Id), teamIdFloat)
	})

	testutil.Run(t, "clears team as admin", func(t *testing.T) {
		teamFactory := factories.NewTeamFactory()
		team, err := teamFactory.Create(map[string]any{"title": "To Clear"})
		require.NoError(t, err)

		projectFactory := factories.NewProjectFactory()
		project, err := projectFactory.Create(map[string]any{
			"name": "With Team",
			"team": team,
		})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		body := map[string]any{
			"team_id": nil,
		}
		path := fmt.Sprintf("/api/projects/%d/set_team", project.Id)
		resp := testutil.ReqJSON(t, suite.Server, "POST", path, body, &token)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Nil(t, result["team_id"])
	})

	testutil.Run(t, "returns 401 without auth", func(t *testing.T) {
		body := map[string]any{
			"team_id": 1,
		}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/projects/1/set_team", body, nil)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
