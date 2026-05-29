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

func TestGetGroups(t *testing.T) {
	testutil.Run(t, "returns empty list when no groups", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/groups", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Empty(t, result)
	})

	testutil.Run(t, "returns all groups", func(t *testing.T) {
		factory := factories.NewGroupFactory()
		group1, err := factory.Create(map[string]any{"name": "Group Alpha"})
		require.NoError(t, err)

		group2, err := factory.Create(map[string]any{"name": "Group Beta"})
		require.NoError(t, err)

		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/groups", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result, 2)

		ids := make([]uint, len(result))
		for i, r := range result {
			ids[i] = uint(r["id"].(float64))
		}

		assert.Contains(t, ids, uint(group1.Id))
		assert.Contains(t, ids, uint(group2.Id))
	})
}

func TestGetGroupById(t *testing.T) {
	testutil.Run(t, "returns existing group", func(t *testing.T) {
		factory := factories.NewGroupFactory()
		group, err := factory.Create(map[string]any{"name": "Test Group"})
		require.NoError(t, err)

		path := fmt.Sprintf("/api/groups/%d", group.Id)
		resp := testutil.ReqJSON(t, suite.Server, "GET", path, nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Equal(t, float64(group.Id), result["id"])
		assert.Equal(t, "Test Group", result["title"])
	})

	testutil.Run(t, "returns 404 for nonexistent group", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/groups/999999", nil, nil)

		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.NotEmpty(t, result["error"])
	})
}
