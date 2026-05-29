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

func TestGetUsers(t *testing.T) {
	testutil.Run(t, "returns empty list when no users", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/users", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Empty(t, result)
	})

	testutil.Run(t, "returns all users", func(t *testing.T) {
		factory := factories.NewUserFactory()
		user1, err := factory.Create(map[string]any{"name": "User Alpha"})
		require.NoError(t, err)

		user2, err := factory.Create(map[string]any{"name": "User Beta"})
		require.NoError(t, err)

		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/users", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result, 2)

		ids := make([]uint, len(result))
		for i, r := range result {
			ids[i] = uint(r["id"].(float64))
		}

		assert.Contains(t, ids, uint(user1.Id))
		assert.Contains(t, ids, uint(user2.Id))
	})

	testutil.Run(t, "filters by search", func(t *testing.T) {
		factory := factories.NewUserFactory()
		_, err := factory.Create(map[string]any{"name": "Searchable First"})
		require.NoError(t, err)

		_, err = factory.Create(map[string]any{"name": "Searchable Second"})
		require.NoError(t, err)

		_, err = factory.Create(map[string]any{"name": "Other User"})
		require.NoError(t, err)

		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/users?search=First", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result, 1)
		assert.Equal(t, "Searchable First", result[0]["name"])
	})

	testutil.Run(t, "filters by team_id returns empty without matching assignees", func(t *testing.T) {
		userFactory := factories.NewUserFactory()
		_, err := userFactory.Create(map[string]any{"name": "Some User"})
		require.NoError(t, err)

		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/users?team_id=999", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Empty(t, result)
	})
}

func TestGetUser(t *testing.T) {
	testutil.Run(t, "returns existing user", func(t *testing.T) {
		factory := factories.NewUserFactory()
		user, err := factory.Create(map[string]any{"name": "Test User"})
		require.NoError(t, err)

		path := fmt.Sprintf("/api/users/%d", user.Id)
		resp := testutil.ReqJSON(t, suite.Server, "GET", path, nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Equal(t, float64(user.Id), result["id"])
		assert.Equal(t, "Test User", result["name"])
	})

	testutil.Run(t, "returns 404 for nonexistent user", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/users/999999", nil, nil)

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestSetUserVisibility(t *testing.T) {
	testutil.Run(t, "sets visibility as admin", func(t *testing.T) {
		factory := factories.NewUserFactory()
		user, err := factory.Create(map[string]any{
			"name":       "Visible User",
			"is_visible": true,
		})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		body := map[string]any{"visible": false}
		path := fmt.Sprintf("/api/users/%d/visibility", user.Id)
		resp := testutil.ReqJSON(t, suite.Server, "POST", path, body, &token)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Equal(t, false, result["is_visible"])
	})

	testutil.Run(t, "returns 401 without auth", func(t *testing.T) {
		body := map[string]any{"visible": false}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/users/1/visibility", body, nil)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestSetUserGroups(t *testing.T) {
	testutil.Run(t, "sets groups as admin", func(t *testing.T) {
		groupFactory := factories.NewGroupFactory()
		group1, err := groupFactory.Create(map[string]any{"name": "Backend"})
		require.NoError(t, err)

		group2, err := groupFactory.Create(map[string]any{"name": "Frontend"})
		require.NoError(t, err)

		userFactory := factories.NewUserFactory()
		user, err := userFactory.Create(map[string]any{"name": "Group User"})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		body := map[string]any{
			"groups": []uint{uint(group1.Id), uint(group2.Id)},
		}
		path := fmt.Sprintf("/api/users/%d/groups", user.Id)
		resp := testutil.ReqJSON(t, suite.Server, "POST", path, body, &token)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		groups := result["groups"].([]any)
		require.Len(t, groups, 2)
	})

	testutil.Run(t, "returns 401 without auth", func(t *testing.T) {
		body := map[string]any{"groups": []uint{1}}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/users/1/groups", body, nil)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestSetUserRole(t *testing.T) {
	testutil.Run(t, "sets role as admin", func(t *testing.T) {
		userFactory := factories.NewUserFactory()
		user, err := userFactory.Create(map[string]any{"name": "Role User"})
		require.NoError(t, err)

		accountFactory := factories.NewAccountFactory()
		_, err = accountFactory.Create(map[string]any{
			"gitlab_id": uint(user.Id),
			"role":      "viewer",
		})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		body := map[string]any{"role": "admin"}
		path := fmt.Sprintf("/api/users/%d/role", user.Id)
		resp := testutil.ReqJSON(t, suite.Server, "POST", path, body, &token)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		account := result["account"].(map[string]any)
		assert.Equal(t, "admin", account["role"])
	})

	testutil.Run(t, "returns 401 without auth", func(t *testing.T) {
		body := map[string]any{"role": "admin"}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/users/1/role", body, nil)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	testutil.Run(t, "returns 403 for non-admin", func(t *testing.T) {
		token := testutil.GetViewerToken(t, suite.Server)

		body := map[string]any{"role": "admin"}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/users/1/role", body, &token)

		require.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}
