package integration

import (
	"main/internal/testutil"
	"main/internal/testutil/factories"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLabels(t *testing.T) {
	testutil.Run(t, "returns empty list when no labels", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/labels", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Empty(t, result)
	})

	testutil.Run(t, "returns unique labels by name", func(t *testing.T) {
		factory := factories.NewLabelFactory()
		labelName := "bug"

		_, err := factory.Create(map[string]any{
			"id":    "uuid-1",
			"name":  labelName,
			"color": "#FF0000",
		})
		require.NoError(t, err)

		_, err = factory.Create(map[string]any{
			"id":    "uuid-2",
			"name":  labelName,
			"color": "#00FF00",
		})
		require.NoError(t, err)

		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/labels", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]any
		testutil.ParseBody(t, resp, &result)

		require.Len(t, result, 1)
		assert.Equal(t, "bug", result[0]["name"])
	})
}

func TestUpdateLabel(t *testing.T) {
	testutil.Run(t, "updates label as admin", func(t *testing.T) {
		factory := factories.NewLabelFactory()
		_, err := factory.Create(map[string]any{
			"id":   "update-label-id",
			"name": "feature",
		})
		require.NoError(t, err)

		token := testutil.GetAdminToken(t, suite.Server)

		body := map[string]any{
			"altName":       "Feature Request",
			"bindingStatus": "backlog",
			"priority":      "high",
		}
		resp := testutil.ReqJSON(t, suite.Server, "PUT", "/api/labels/feature", body, &token)

		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	testutil.Run(t, "returns 401 without auth", func(t *testing.T) {
		body := map[string]any{
			"altName": "Test",
		}
		resp := testutil.ReqJSON(t, suite.Server, "PUT", "/api/labels/test", body, nil)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	testutil.Run(t, "returns 403 for non-admin", func(t *testing.T) {
		token := testutil.GetViewerToken(t, suite.Server)

		body := map[string]any{
			"altName": "Test",
		}
		resp := testutil.ReqJSON(t, suite.Server, "PUT", "/api/labels/test", body, &token)

		require.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}
