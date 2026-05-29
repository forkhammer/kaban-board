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

func TestGetEpics(t *testing.T) {
	testutil.Run(t, "returns empty list when no epics", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/epic", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Empty(t, result)
	})

	testutil.Run(t, "returns all epics", func(t *testing.T) {
		projectFactory := factories.NewProjectFactory()
		project, err := projectFactory.Create(nil)
		require.NoError(t, err)

		epicFactory := factories.NewEpicFactory()
		epic1, err := epicFactory.Create(map[string]any{
			"project": *project,
			"title":   "Epic One",
		})
		require.NoError(t, err)

		epic2, err := epicFactory.Create(map[string]any{
			"project": *project,
			"title":   "Epic Two",
		})
		require.NoError(t, err)

		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/epic", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result, 2)

		ids := make([]uint, len(result))
		for i, r := range result {
			ids[i] = uint(r["id"].(float64))
		}

		assert.Contains(t, ids, uint(epic1.Id))
		assert.Contains(t, ids, uint(epic2.Id))
	})

	testutil.Run(t, "filters by project", func(t *testing.T) {
		projectFactory := factories.NewProjectFactory()
		projectA, err := projectFactory.Create(map[string]any{"name": "Project A"})
		require.NoError(t, err)

		projectB, err := projectFactory.Create(map[string]any{"name": "Project B"})
		require.NoError(t, err)

		epicFactory := factories.NewEpicFactory()
		_, err = epicFactory.Create(map[string]any{
			"project": *projectA,
			"title":   "Alpha Epic A",
		})
		require.NoError(t, err)

		_, err = epicFactory.Create(map[string]any{
			"project": *projectA,
			"title":   "Beta Epic A",
		})
		require.NoError(t, err)

		_, err = epicFactory.Create(map[string]any{
			"project": *projectB,
			"title":   "Epic B",
		})
		require.NoError(t, err)

		path := fmt.Sprintf("/api/epic?project=%d", projectA.Id)
		resp := testutil.ReqJSON(t, suite.Server, "GET", path, nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result, 2)

		titles := make([]string, len(result))
		for i, r := range result {
			titles[i] = r["title"].(string)
		}

		assert.Contains(t, titles, "Alpha Epic A")
		assert.Contains(t, titles, "Beta Epic A")
		assert.NotContains(t, titles, "Epic B")
	})

	testutil.Run(t, "filters by search", func(t *testing.T) {
		projectFactory := factories.NewProjectFactory()
		project, err := projectFactory.Create(nil)
		require.NoError(t, err)

		epicFactory := factories.NewEpicFactory()
		_, err = epicFactory.Create(map[string]any{
			"project": *project,
			"title":   "Searchable First",
		})
		require.NoError(t, err)

		_, err = epicFactory.Create(map[string]any{
			"project": *project,
			"title":   "Searchable Second",
		})
		require.NoError(t, err)

		_, err = epicFactory.Create(map[string]any{
			"project": *project,
			"title":   "Other Epic",
		})
		require.NoError(t, err)

		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/epic?search=First", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result, 1)
		assert.Equal(t, "Searchable First", result[0]["title"])
	})

	testutil.Run(t, "filters by project and search combined", func(t *testing.T) {
		projectFactory := factories.NewProjectFactory()
		projectA, err := projectFactory.Create(map[string]any{"name": "Project A"})
		require.NoError(t, err)

		projectB, err := projectFactory.Create(map[string]any{"name": "Project B"})
		require.NoError(t, err)

		epicFactory := factories.NewEpicFactory()
		_, err = epicFactory.Create(map[string]any{
			"project": *projectA,
			"title":   "Combined Test",
		})
		require.NoError(t, err)

		_, err = epicFactory.Create(map[string]any{
			"project": *projectA,
			"title":   "Other A",
		})
		require.NoError(t, err)

		_, err = epicFactory.Create(map[string]any{
			"project": *projectB,
			"title":   "Combined Test",
		})
		require.NoError(t, err)

		path := fmt.Sprintf("/api/epic?project=%d&search=Combined", projectA.Id)
		resp := testutil.ReqJSON(t, suite.Server, "GET", path, nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result, 1)
		assert.Equal(t, "Combined Test", result[0]["title"])
	})
}

func TestGetEpic(t *testing.T) {
	testutil.Run(t, "returns existing epic", func(t *testing.T) {
		projectFactory := factories.NewProjectFactory()
		project, err := projectFactory.Create(map[string]any{"name": "Test Project"})
		require.NoError(t, err)

		epicFactory := factories.NewEpicFactory()
		epic, err := epicFactory.Create(map[string]any{
			"project": *project,
			"title":   "My Epic",
			"iid":     "42",
		})
		require.NoError(t, err)

		path := fmt.Sprintf("/api/epic/%d", epic.Id)
		resp := testutil.ReqJSON(t, suite.Server, "GET", path, nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Equal(t, float64(epic.Id), result["id"])
		assert.Equal(t, "42", result["iid"])
		assert.Equal(t, "My Epic", result["title"])
		assert.Equal(t, "Test Project", result["project"])
	})

	testutil.Run(t, "returns 404 for nonexistent epic", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/epic/999999", nil, nil)

		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.NotEmpty(t, result["error"])
	})

	testutil.Run(t, "returns 400 for invalid id", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/epic/abc", nil, nil)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.NotEmpty(t, result["error"])
	})
}
