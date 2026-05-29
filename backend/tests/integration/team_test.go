package integration

import (
	"fmt"
	domain_models "main/internal/domain/models"
	"main/internal/testutil"
	"main/internal/testutil/factories"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTeams(t *testing.T) {
	testutil.Run(t, "returns empty list when no teams", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/teams", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Empty(t, result)
	})

	testutil.Run(t, "returns teams with groups", func(t *testing.T) {
		groupFactory := factories.NewGroupFactory()
		group1, err := groupFactory.Create(map[string]any{"name": "Backend"})
		require.NoError(t, err)

		group2, err := groupFactory.Create(map[string]any{"name": "Frontend"})
		require.NoError(t, err)

		teamFactory := factories.NewTeamFactory()
		team, err := teamFactory.Create(map[string]any{
			"title":  "Dev Team",
			"groups": []domain_models.Group{*group1, *group2},
		})
		require.NoError(t, err)

		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/teams", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result, 1)
		assert.Equal(t, float64(team.Id), result[0]["id"])
		assert.Equal(t, "Dev Team", result[0]["title"])

		groups := result[0]["groups"].([]any)
		require.Len(t, groups, 2)
	})
}

func TestGetTeamById(t *testing.T) {
	testutil.Run(t, "returns existing team", func(t *testing.T) {
		teamFactory := factories.NewTeamFactory()
		team, err := teamFactory.Create(map[string]any{"title": "Test Team"})
		require.NoError(t, err)

		path := fmt.Sprintf("/api/teams/%d", team.Id)
		resp := testutil.ReqJSON(t, suite.Server, "GET", path, nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Equal(t, float64(team.Id), result["id"])
		assert.Equal(t, "Test Team", result["title"])
	})

	testutil.Run(t, "returns 404 for nonexistent team", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/teams/999999", nil, nil)

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestCreateTeam(t *testing.T) {
	testutil.Run(t, "creates team as admin", func(t *testing.T) {
		token := testutil.GetAdminToken(t, suite.Server)

		body := map[string]any{
			"title":  "New Team",
			"groups": []int{},
		}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/teams", body, &token)

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Equal(t, "New Team", result["title"])
	})

	testutil.Run(t, "returns 401 without auth", func(t *testing.T) {
		body := map[string]any{
			"title":  "New Team",
			"groups": []int{},
		}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/teams", body, nil)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestUpdateTeam(t *testing.T) {
	testutil.Run(t, "updates team as admin", func(t *testing.T) {
		teamFactory := factories.NewTeamFactory()
		team, err := teamFactory.Create(map[string]any{"title": "Old Title"})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		body := map[string]any{
			"title":  "Updated Title",
			"groups": []int{},
		}
		path := fmt.Sprintf("/api/teams/%d", team.Id)
		resp := testutil.ReqJSON(t, suite.Server, "PUT", path, body, &token)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Equal(t, "Updated Title", result["title"])
	})

	testutil.Run(t, "returns 401 without auth", func(t *testing.T) {
		body := map[string]any{
			"title":  "Updated Title",
			"groups": []int{},
		}
		resp := testutil.ReqJSON(t, suite.Server, "PUT", "/api/teams/1", body, nil)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestDeleteTeam(t *testing.T) {
	testutil.Run(t, "deletes team as admin", func(t *testing.T) {
		teamFactory := factories.NewTeamFactory()
		team, err := teamFactory.Create(map[string]any{"title": "To Delete"})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		path := fmt.Sprintf("/api/teams/%d", team.Id)
		resp := testutil.ReqJSON(t, suite.Server, "DELETE", path, nil, &token)

		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	testutil.Run(t, "returns 401 without auth", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "DELETE", "/api/teams/1", nil, nil)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
