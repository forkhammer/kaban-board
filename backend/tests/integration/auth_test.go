package integration

import (
	"main/internal/testutil"
	"main/internal/testutil/factories"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	testutil.Run(t, "valid registration", func(t *testing.T) {
		body := map[string]string{
			"username": "testuser",
			"password": "testpass123",
		}

		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/account/register", body, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Equal(t, "testuser", result["username"])
	})

	testutil.Run(t, "missing username", func(t *testing.T) {
		body := map[string]string{
			"password": "testpass123",
		}

		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/account/register", body, nil)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	testutil.Run(t, "missing password", func(t *testing.T) {
		body := map[string]string{
			"username": "testuser2",
		}

		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/account/register", body, nil)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestLogin(t *testing.T) {
	testutil.Run(t, "valid login", func(t *testing.T) {
		factory := factories.NewAccountFactory()
		account, _ := factory.Create(map[string]any{
			"password":  "loginpass123",
			"is_active": true,
		})

		loginBody := map[string]string{
			"username": account.Username,
			"password": "loginpass123",
		}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/account/login", loginBody, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		require.NotEmpty(t, result["token"])
	})

	testutil.Run(t, "wrong password", func(t *testing.T) {
		factory := factories.NewAccountFactory()
		account, _ := factory.Create(map[string]any{
			"password":  "correctpass",
			"is_active": true,
		})

		loginBody := map[string]string{
			"username": account.Username,
			"password": "wrongpass",
		}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/account/login", loginBody, nil)

		require.NotEqual(t, http.StatusOK, resp.StatusCode)
	})

	testutil.Run(t, "nonexistent user", func(t *testing.T) {
		loginBody := map[string]string{
			"username": "nonexistent",
			"password": "somepass",
		}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/account/login", loginBody, nil)

		require.NotEqual(t, http.StatusOK, resp.StatusCode)
	})
}

func TestGetActiveUser(t *testing.T) {
	testutil.Run(t, "with valid token", func(t *testing.T) {
		factory := factories.NewAccountFactory()
		account, _ := factory.Create(map[string]any{
			"password":  "activepass123",
			"is_active": true,
		})

		loginBody := map[string]string{
			"username": account.Username,
			"password": "activepass123",
		}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/account/login", loginBody, nil)
		var loginResult map[string]any
		testutil.ParseBody(t, resp, &loginResult)

		token := loginResult["token"].(string)

		resp = testutil.ReqJSON(t, suite.Server, "GET", "/api/account/user", nil, &token)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		user := result["user"].(map[string]any)
		assert.Equal(t, account.Username, user["username"])
	})

	testutil.Run(t, "without token", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/account/user", nil, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Nil(t, result["user"])
	})
}

func TestProtectedRoute(t *testing.T) {
	testutil.Run(t, "without auth returns 401", func(t *testing.T) {
		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/account/online", nil, nil)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
