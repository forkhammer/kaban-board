package integration

import (
	"fmt"
	"main/internal/domain/models"
	"main/internal/testutil"
	"main/internal/testutil/factories"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func uintPtr(v uint) *uint { return &v }

type workloadFixture struct {
	team    *models.Team
	sprint  *models.Sprint
	project *models.Project
	issue   *models.Issue
	user1   *models.User
}

func createWorkloadFixture(t *testing.T) workloadFixture {
	t.Helper()

	team, err := factories.NewTeamFactory().Create(nil)
	require.NoError(t, err)

	sprint, err := factories.NewSprintFactory().Create(map[string]any{
		"team":           *team,
		"hours_per_user": uint(40),
	})
	require.NoError(t, err)

	project, err := factories.NewProjectFactory().Create(map[string]any{"team": team})
	require.NoError(t, err)

	issue, err := factories.NewIssueFactory().Create(map[string]any{"project": *project})
	require.NoError(t, err)

	user1, err := factories.NewUserFactory().Create(map[string]any{"name": "Workload User One"})
	require.NoError(t, err)

	_, err = factories.NewUserFactory().Create(map[string]any{"name": "Workload User Two"})
	require.NoError(t, err)

	_, err = factories.NewIssueBindingFactory().Create(map[string]any{
		"sprint":       sprint,
		"issue":        issue,
		"assignee":     user1,
		"estimate_dev": uintPtr(5),
		"estimate_qa":  uintPtr(3),
	})
	require.NoError(t, err)

	return workloadFixture{
		team:    team,
		sprint:  sprint,
		project: project,
		issue:   issue,
		user1:   user1,
	}
}

func TestSprintUsersWorkload(t *testing.T) {
	testutil.Run(t, "returns planned and capacity per assigned user", func(t *testing.T) {
		fx := createWorkloadFixture(t)

		path := fmt.Sprintf("/api/reports/sprint/%d/users-workload", fx.sprint.Id)
		resp := testutil.ReqJSON(t, suite.Server, "GET", path, nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result struct {
			Users []map[string]any `json:"users"`
		}
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result.Users, 1)
		assert.Equal(t, float64(fx.user1.Id), result.Users[0]["user_id"])
		assert.Equal(t, float64(8), result.Users[0]["planned"])
		assert.Equal(t, float64(40), result.Users[0]["capacity"])
	})

	testutil.Run(t, "sums planned across multiple bindings", func(t *testing.T) {
		fx := createWorkloadFixture(t)

		issue2, err := factories.NewIssueFactory().Create(map[string]any{"project": *fx.project})
		require.NoError(t, err)

		_, err = factories.NewIssueBindingFactory().Create(map[string]any{
			"sprint":       fx.sprint,
			"issue":        issue2,
			"assignee":     fx.user1,
			"estimate_dev": uintPtr(2),
			"estimate_qa":  uintPtr(0),
		})
		require.NoError(t, err)

		path := fmt.Sprintf("/api/reports/sprint/%d/users-workload", fx.sprint.Id)
		resp := testutil.ReqJSON(t, suite.Server, "GET", path, nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result struct {
			Users []map[string]any `json:"users"`
		}
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result.Users, 1)
		assert.Equal(t, float64(10), result.Users[0]["planned"])
	})

	testutil.Run(t, "uses per-user settings capacity when present", func(t *testing.T) {
		fx := createWorkloadFixture(t)

		token := testutil.GetAdminToken(t, suite.Server)

		body := map[string]any{
			"user_id":        fx.user1.Id,
			"hours_per_user": 30,
		}
		settingsPath := fmt.Sprintf("/api/sprint/%d/user-settings", fx.sprint.Id)
		resp := testutil.ReqJSON(t, suite.Server, "POST", settingsPath, body, &token)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		path := fmt.Sprintf("/api/reports/sprint/%d/users-workload", fx.sprint.Id)
		resp = testutil.ReqJSON(t, suite.Server, "GET", path, nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result struct {
			Users []map[string]any `json:"users"`
		}
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result.Users, 1)
		assert.Equal(t, float64(30), result.Users[0]["capacity"])
	})

	testutil.Run(t, "returns empty users array for sprint with no bindings", func(t *testing.T) {
		team, err := factories.NewTeamFactory().Create(nil)
		require.NoError(t, err)

		sprint, err := factories.NewSprintFactory().Create(map[string]any{"team": *team})
		require.NoError(t, err)

		path := fmt.Sprintf("/api/reports/sprint/%d/users-workload", sprint.Id)
		resp := testutil.ReqJSON(t, suite.Server, "GET", path, nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result struct {
			Users []map[string]any `json:"users"`
		}
		testutil.ParseBody(t, resp, &result)

		assert.NotNil(t, result.Users)
		assert.Empty(t, result.Users)
	})

	testutil.Run(t, "returns 404 for nonexistent sprint", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/reports/sprint/999999/users-workload", nil, nil)

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	testutil.Run(t, "returns 400 for invalid sprint id", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/reports/sprint/abc/users-workload", nil, nil)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
