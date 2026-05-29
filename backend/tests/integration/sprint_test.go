package integration

import (
	"fmt"
	"main/internal/domain/models"
	"main/internal/testutil"
	"main/internal/testutil/factories"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSprints(t *testing.T) {
	testutil.Run(t, "returns empty list when no sprints", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/sprint", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Empty(t, result)
	})

	testutil.Run(t, "returns all sprints", func(t *testing.T) {
		teamFactory := factories.NewTeamFactory()
		team, err := teamFactory.Create(map[string]any{"title": "Sprint Team"})
		require.NoError(t, err)

		sprintFactory := factories.NewSprintFactory()
		sprint1, err := sprintFactory.Create(map[string]any{
			"title":        "Sprint One",
			"team":         *team,
			"start_date":   time.Now().AddDate(0, 0, -7),
			"end_date":     time.Now().AddDate(0, 0, 7),
			"hours_per_user": 40,
		})
		require.NoError(t, err)

		sprint2, err := sprintFactory.Create(map[string]any{
			"title":        "Sprint Two",
			"team":         *team,
			"start_date":   time.Now().AddDate(0, 0, -14),
			"end_date":     time.Now().AddDate(0, 0, -7),
			"hours_per_user": 40,
		})
		require.NoError(t, err)

		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/sprint", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result, 2)

		ids := make([]int, len(result))
		for i, r := range result {
			ids[i] = int(r["id"].(float64))
		}

		assert.Contains(t, ids, int(sprint1.Id))
		assert.Contains(t, ids, int(sprint2.Id))
	})

	testutil.Run(t, "filters by team", func(t *testing.T) {
		teamFactory := factories.NewTeamFactory()
		teamA, err := teamFactory.Create(map[string]any{"title": "Team A"})
		require.NoError(t, err)

		teamB, err := teamFactory.Create(map[string]any{"title": "Team B"})
		require.NoError(t, err)

		sprintFactory := factories.NewSprintFactory()
		_, err = sprintFactory.Create(map[string]any{
			"title":      "Sprint A",
			"team":       *teamA,
			"start_date": time.Now().AddDate(0, 0, -7),
			"end_date":   time.Now().AddDate(0, 0, 7),
		})
		require.NoError(t, err)

		_, err = sprintFactory.Create(map[string]any{
			"title":      "Sprint B",
			"team":       *teamB,
			"start_date": time.Now().AddDate(0, 0, -7),
			"end_date":   time.Now().AddDate(0, 0, 7),
		})
		require.NoError(t, err)

		path := fmt.Sprintf("/api/sprint?team=%d", int(teamA.Id))
		resp := testutil.ReqJSON(t, suite.Server, "GET", path, nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result, 1)
		assert.Equal(t, "Sprint A", result[0]["title"])
	})

	testutil.Run(t, "filters by quarter", func(t *testing.T) {
		teamFactory := factories.NewTeamFactory()
		team, err := teamFactory.Create(map[string]any{"title": "Team Q"})
		require.NoError(t, err)

		sprintFactory := factories.NewSprintFactory()
		_, err = sprintFactory.Create(map[string]any{
			"title":      "Q1 Sprint",
			"team":       *team,
			"start_date": time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			"end_date":   time.Date(2024, 1, 28, 0, 0, 0, 0, time.UTC),
		})
		require.NoError(t, err)

		_, err = sprintFactory.Create(map[string]any{
			"title":      "Q2 Sprint",
			"team":       *team,
			"start_date": time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC),
			"end_date":   time.Date(2024, 4, 28, 0, 0, 0, 0, time.UTC),
		})
		require.NoError(t, err)

		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/sprint?quarter=2024-1", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result, 1)
		assert.Equal(t, "Q1 Sprint", result[0]["title"])
	})
}

func TestGetSprint(t *testing.T) {
	testutil.Run(t, "returns existing sprint", func(t *testing.T) {
		teamFactory := factories.NewTeamFactory()
		team, err := teamFactory.Create(map[string]any{"title": "Test Team"})
		require.NoError(t, err)

		sprintFactory := factories.NewSprintFactory()
		sprint, err := sprintFactory.Create(map[string]any{
			"title":        "My Sprint",
			"team":         *team,
			"start_date":   time.Now().AddDate(0, 0, -7),
			"end_date":     time.Now().AddDate(0, 0, 7),
			"hours_per_user": 40,
		})
		require.NoError(t, err)

		path := fmt.Sprintf("/api/sprint/%d", sprint.Id)
		resp := testutil.ReqJSON(t, suite.Server, "GET", path, nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Equal(t, float64(sprint.Id), result["id"])
		assert.Equal(t, "My Sprint", result["title"])
	})

	testutil.Run(t, "returns 404 for nonexistent sprint", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/sprint/999999", nil, nil)

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestGetQuarters(t *testing.T) {
	testutil.Run(t, "returns quarters", func(t *testing.T) {
		teamFactory := factories.NewTeamFactory()
		team, err := teamFactory.Create(map[string]any{"title": "Q Team"})
		require.NoError(t, err)

		sprintFactory := factories.NewSprintFactory()
		_, err = sprintFactory.Create(map[string]any{
			"title":      "Q1 Sprint",
			"team":       *team,
			"start_date": time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			"end_date":   time.Date(2024, 1, 28, 0, 0, 0, 0, time.UTC),
		})
		require.NoError(t, err)

		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/sprint/quarters", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		require.NotEmpty(t, result)
	})
}

func TestCreateSprint(t *testing.T) {
	testutil.Run(t, "creates sprint as admin", func(t *testing.T) {
		teamFactory := factories.NewTeamFactory()
		team, err := teamFactory.Create(map[string]any{"title": "Create Team"})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		body := map[string]any{
			"title":          "New Sprint",
			"start_date":     time.Now().AddDate(0, 0, -7).Format("2006-01-02"),
			"end_date":       time.Now().AddDate(0, 0, 7).Format("2006-01-02"),
			"hours_per_user": 40,
			"team_id":        team.Id,
		}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/sprint", body, &token)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Equal(t, "New Sprint", result["title"])
		assert.Equal(t, float64(team.Id), result["team_id"])
	})

	testutil.Run(t, "returns 401 without auth", func(t *testing.T) {
		body := map[string]any{
			"title":          "Sprint",
			"start_date":     time.Now().Format("2006-01-02"),
			"end_date":       time.Now().AddDate(0, 0, 7).Format("2006-01-02"),
			"hours_per_user": 40,
			"team_id":        1,
		}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/sprint", body, nil)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	testutil.Run(t, "returns 400 for invalid dates", func(t *testing.T) {
		teamFactory := factories.NewTeamFactory()
		team, err := teamFactory.Create(map[string]any{"title": "Bad Date Team"})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		body := map[string]any{
			"title":          "Bad Sprint",
			"start_date":     "invalid-date",
			"end_date":       "invalid-date",
			"hours_per_user": 40,
			"team_id":        team.Id,
		}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/sprint", body, &token)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestUpdateSprint(t *testing.T) {
	testutil.Run(t, "updates sprint as admin", func(t *testing.T) {
		teamFactory := factories.NewTeamFactory()
		team, err := teamFactory.Create(map[string]any{"title": "Update Team"})
		require.NoError(t, err)

		sprintFactory := factories.NewSprintFactory()
		sprint, err := sprintFactory.Create(map[string]any{
			"title":        "Old Sprint",
			"team":         *team,
			"start_date":   time.Now().AddDate(0, 0, -7),
			"end_date":     time.Now().AddDate(0, 0, 7),
		})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		body := map[string]any{
			"title":          "Updated Sprint",
			"start_date":     time.Now().AddDate(0, 0, -7).Format("2006-01-02"),
			"end_date":       time.Now().AddDate(0, 0, 7).Format("2006-01-02"),
			"hours_per_user": 40,
			"team_id":        team.Id,
		}
		path := fmt.Sprintf("/api/sprint/%d", sprint.Id)
		resp := testutil.ReqJSON(t, suite.Server, "PUT", path, body, &token)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Equal(t, "Updated Sprint", result["title"])
	})

	testutil.Run(t, "returns 401 without auth", func(t *testing.T) {
		body := map[string]any{
			"title":          "Sprint",
			"start_date":     time.Now().Format("2006-01-02"),
			"end_date":       time.Now().AddDate(0, 0, 7).Format("2006-01-02"),
			"hours_per_user": 40,
			"team_id":        1,
		}
		resp := testutil.ReqJSON(t, suite.Server, "PUT", "/api/sprint/1", body, nil)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestDeleteSprint(t *testing.T) {
	testutil.Run(t, "deletes sprint as admin", func(t *testing.T) {
		teamFactory := factories.NewTeamFactory()
		team, err := teamFactory.Create(map[string]any{"title": "Delete Team"})
		require.NoError(t, err)

		sprintFactory := factories.NewSprintFactory()
		sprint, err := sprintFactory.Create(map[string]any{
			"title":        "To Delete",
			"team":         *team,
			"start_date":   time.Now().AddDate(0, 0, -7),
			"end_date":     time.Now().AddDate(0, 0, 7),
			"status":       models.SprintStatusWaiting,
		})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		path := fmt.Sprintf("/api/sprint/%d", sprint.Id)
		resp := testutil.ReqJSON(t, suite.Server, "DELETE", path, nil, &token)

		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	testutil.Run(t, "returns 401 without auth", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "DELETE", "/api/sprint/1", nil, nil)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestCompleteSprint(t *testing.T) {
	testutil.Run(t, "completes running sprint as admin", func(t *testing.T) {
		teamFactory := factories.NewTeamFactory()
		team, err := teamFactory.Create(map[string]any{"title": "Complete Team"})
		require.NoError(t, err)

		sprintFactory := factories.NewSprintFactory()
		sprint, err := sprintFactory.Create(map[string]any{
			"title":        "To Complete",
			"team":         *team,
			"start_date":   time.Now().AddDate(0, 0, -7),
			"end_date":     time.Now().AddDate(0, 0, 7),
			"status":       models.SprintStatusRunning,
		})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		path := fmt.Sprintf("/api/sprint/%d/complete", sprint.Id)
		resp := testutil.ReqJSON(t, suite.Server, "POST", path, nil, &token)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Equal(t, "completed", result["status"])
	})

	testutil.Run(t, "returns 401 without auth", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/sprint/1/complete", nil, nil)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestRunSprint(t *testing.T) {
	testutil.Run(t, "runs waiting sprint as admin", func(t *testing.T) {
		teamFactory := factories.NewTeamFactory()
		team, err := teamFactory.Create(map[string]any{"title": "Run Team"})
		require.NoError(t, err)

		sprintFactory := factories.NewSprintFactory()
		sprint, err := sprintFactory.Create(map[string]any{
			"title":        "To Run",
			"team":         *team,
			"start_date":   time.Now().AddDate(0, 0, -1),
			"end_date":     time.Now().AddDate(0, 0, 7),
			"status":       models.SprintStatusWaiting,
		})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		path := fmt.Sprintf("/api/sprint/%d/run", sprint.Id)
		resp := testutil.ReqJSON(t, suite.Server, "POST", path, nil, &token)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Equal(t, "running", result["status"])
	})

	testutil.Run(t, "returns 401 without auth", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/sprint/1/run", nil, nil)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestSprintUserSettings(t *testing.T) {
	testutil.Run(t, "creates and gets user settings", func(t *testing.T) {
		teamFactory := factories.NewTeamFactory()
		team, err := teamFactory.Create(map[string]any{"title": "Settings Team"})
		require.NoError(t, err)

		sprintFactory := factories.NewSprintFactory()
		sprint, err := sprintFactory.Create(map[string]any{
			"title":      "Settings Sprint",
			"team":       *team,
			"start_date": time.Now().AddDate(0, 0, -7),
			"end_date":   time.Now().AddDate(0, 0, 7),
		})
		require.NoError(t, err)

		userFactory := factories.NewUserFactory()
		user, err := userFactory.Create(map[string]any{"name": "Settings User"})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		body := map[string]any{
			"user_id":        user.Id,
			"hours_per_user": 30,
		}
		path := fmt.Sprintf("/api/sprint/%d/user-settings", sprint.Id)
		resp := testutil.ReqJSON(t, suite.Server, "POST", path, body, &token)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var createResult map[string]any
		testutil.ParseBody(t, resp, &createResult)

		assert.Equal(t, float64(user.Id), createResult["user_id"])
		assert.Equal(t, float64(30), createResult["hours_per_user"])

		resp = testutil.ReqJSON(t, suite.Server, "GET", path, nil, &token)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var getResult []map[string]any
		testutil.ParseBody(t, resp, &getResult)

		require.Len(t, getResult, 1)
		assert.Equal(t, float64(user.Id), getResult[0]["user_id"])
	})

	testutil.Run(t, "deletes user settings", func(t *testing.T) {
		teamFactory := factories.NewTeamFactory()
		team, err := teamFactory.Create(map[string]any{"title": "Delete Settings Team"})
		require.NoError(t, err)

		sprintFactory := factories.NewSprintFactory()
		sprint, err := sprintFactory.Create(map[string]any{
			"title":      "Delete Settings Sprint",
			"team":       *team,
			"start_date": time.Now().AddDate(0, 0, -7),
			"end_date":   time.Now().AddDate(0, 0, 7),
		})
		require.NoError(t, err)

		userFactory := factories.NewUserFactory()
		user, err := userFactory.Create(map[string]any{"name": "Delete User"})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		body := map[string]any{
			"user_id":        user.Id,
			"hours_per_user": 25,
		}
		path := fmt.Sprintf("/api/sprint/%d/user-settings", sprint.Id)
		resp := testutil.ReqJSON(t, suite.Server, "POST", path, body, &token)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		delPath := fmt.Sprintf("/api/sprint/%d/user-settings/%d", sprint.Id, user.Id)
		resp = testutil.ReqJSON(t, suite.Server, "DELETE", delPath, nil, &token)

		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	testutil.Run(t, "returns 401 without auth for user-settings", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/sprint/1/user-settings", nil, nil)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
