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

func TestGetReleases(t *testing.T) {
	testutil.Run(t, "returns empty list when no releases", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/release", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Empty(t, result)
	})

	testutil.Run(t, "returns all releases", func(t *testing.T) {
		projectFactory := factories.NewProjectFactory()
		project, err := projectFactory.Create(nil)
		require.NoError(t, err)

		releaseFactory := factories.NewReleaseFactory()
		release1, err := releaseFactory.Create(map[string]any{
			"project": *project,
			"title":   "v1.0.0",
		})
		require.NoError(t, err)

		release2, err := releaseFactory.Create(map[string]any{
			"project": *project,
			"title":   "v2.0.0",
		})
		require.NoError(t, err)

		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/release", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result, 2)

		ids := make([]uint, len(result))
		for i, r := range result {
			ids[i] = uint(r["id"].(float64))
		}

		assert.Contains(t, ids, uint(release1.Id))
		assert.Contains(t, ids, uint(release2.Id))
	})

	testutil.Run(t, "filters by project", func(t *testing.T) {
		projectFactory := factories.NewProjectFactory()
		projectA, err := projectFactory.Create(map[string]any{"name": "Project A"})
		require.NoError(t, err)

		projectB, err := projectFactory.Create(map[string]any{"name": "Project B"})
		require.NoError(t, err)

		releaseFactory := factories.NewReleaseFactory()
		_, err = releaseFactory.Create(map[string]any{
			"project": *projectA,
			"title":   "Alpha Release A",
		})
		require.NoError(t, err)

		_, err = releaseFactory.Create(map[string]any{
			"project": *projectA,
			"title":   "Beta Release A",
		})
		require.NoError(t, err)

		_, err = releaseFactory.Create(map[string]any{
			"project": *projectB,
			"title":   "Release B",
		})
		require.NoError(t, err)

		path := fmt.Sprintf("/api/release?project=%d", projectA.Id)
		resp := testutil.ReqJSON(t, suite.Server, "GET", path, nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result, 2)

		titles := make([]string, len(result))
		for i, r := range result {
			titles[i] = r["title"].(string)
		}

		assert.Contains(t, titles, "Alpha Release A")
		assert.Contains(t, titles, "Beta Release A")
		assert.NotContains(t, titles, "Release B")
	})

	testutil.Run(t, "filters by search", func(t *testing.T) {
		projectFactory := factories.NewProjectFactory()
		project, err := projectFactory.Create(nil)
		require.NoError(t, err)

		releaseFactory := factories.NewReleaseFactory()
		_, err = releaseFactory.Create(map[string]any{
			"project": *project,
			"title":   "v1.0.0",
		})
		require.NoError(t, err)

		_, err = releaseFactory.Create(map[string]any{
			"project": *project,
			"title":   "v2.0.0",
		})
		require.NoError(t, err)

		_, err = releaseFactory.Create(map[string]any{
			"project": *project,
			"title":   "alpha",
		})
		require.NoError(t, err)

		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/release?search=1.0", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result, 1)
		assert.Equal(t, "v1.0.0", result[0]["title"])
	})

	testutil.Run(t, "filters by project and search combined", func(t *testing.T) {
		projectFactory := factories.NewProjectFactory()
		projectA, err := projectFactory.Create(map[string]any{"name": "Project A"})
		require.NoError(t, err)

		projectB, err := projectFactory.Create(map[string]any{"name": "Project B"})
		require.NoError(t, err)

		releaseFactory := factories.NewReleaseFactory()
		_, err = releaseFactory.Create(map[string]any{
			"project": *projectA,
			"title":   "v1.0.0",
		})
		require.NoError(t, err)

		_, err = releaseFactory.Create(map[string]any{
			"project": *projectA,
			"title":   "v2.0.0",
		})
		require.NoError(t, err)

		_, err = releaseFactory.Create(map[string]any{
			"project": *projectB,
			"title":   "v1.0.0",
		})
		require.NoError(t, err)

		path := fmt.Sprintf("/api/release?project=%d&search=1.0", projectA.Id)
		resp := testutil.ReqJSON(t, suite.Server, "GET", path, nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result, 1)
		assert.Equal(t, "v1.0.0", result[0]["title"])
	})
}

func TestGetRelease(t *testing.T) {
	testutil.Run(t, "returns existing release", func(t *testing.T) {
		projectFactory := factories.NewProjectFactory()
		project, err := projectFactory.Create(nil)
		require.NoError(t, err)

		releaseFactory := factories.NewReleaseFactory()
		release, err := releaseFactory.Create(map[string]any{
			"project":  *project,
			"title":    "v1.0.0",
			"web_path": "/releases/v1.0.0",
		})
		require.NoError(t, err)

		path := fmt.Sprintf("/api/release/%d", release.Id)
		resp := testutil.ReqJSON(t, suite.Server, "GET", path, nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Equal(t, float64(release.Id), result["id"])
		assert.Equal(t, "v1.0.0", result["title"])
		assert.Equal(t, "/releases/v1.0.0", result["webPath"])
	})

	testutil.Run(t, "returns 404 for nonexistent release", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/release/999999", nil, nil)

		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.NotEmpty(t, result["error"])
	})

	testutil.Run(t, "returns 400 for invalid id", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/release/abc", nil, nil)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.NotEmpty(t, result["error"])
	})
}
