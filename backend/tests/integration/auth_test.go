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
	t.Run("valid registration", func(t *testing.T) {
		t.Cleanup(func() { testutil.CleanupDatabase(t) })

		body := map[string]string{
			"username": "testuser",
			"password": "testpass123",
		}

		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/account/register", body, nil)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Equal(t, "testuser", result["username"])
	})

	t.Run("missing username", func(t *testing.T) {
		t.Cleanup(func() { testutil.CleanupDatabase(t) })

		body := map[string]string{
			"password": "testpass123",
		}

		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/account/register", body, nil)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("missing password", func(t *testing.T) {
		t.Cleanup(func() { testutil.CleanupDatabase(t) })

		body := map[string]string{
			"username": "testuser2",
		}

		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/account/register", body, nil)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestLogin(t *testing.T) {
	t.Run("valid login", func(t *testing.T) {
		t.Cleanup(func() { testutil.CleanupDatabase(t) })

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
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		require.NotEmpty(t, result["token"])
	})

	t.Run("wrong password", func(t *testing.T) {
		t.Cleanup(func() { testutil.CleanupDatabase(t) })

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
		defer resp.Body.Close()

		require.NotEqual(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("nonexistent user", func(t *testing.T) {
		t.Cleanup(func() { testutil.CleanupDatabase(t) })

		loginBody := map[string]string{
			"username": "nonexistent",
			"password": "somepass",
		}
		resp := testutil.ReqJSON(t, suite.Server, "POST", "/api/account/login", loginBody, nil)
		defer resp.Body.Close()

		require.NotEqual(t, http.StatusOK, resp.StatusCode)
	})
}

func TestGetActiveUser(t *testing.T) {
	t.Run("with valid token", func(t *testing.T) {
		t.Cleanup(func() { testutil.CleanupDatabase(t) })

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
		resp.Body.Close()

		token := loginResult["token"].(string)

		resp = testutil.ReqJSON(t, suite.Server, "GET", "/api/account/user", nil, &token)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		user := result["user"].(map[string]any)
		assert.Equal(t, account.Username, user["username"])
	})

	t.Run("without token", func(t *testing.T) {
		t.Cleanup(func() { testutil.CleanupDatabase(t) })

		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/account/user", nil, nil)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		assert.Nil(t, result["user"])
	})
}

func TestProtectedRoute(t *testing.T) {
	t.Run("without auth returns 401", func(t *testing.T) {
		t.Cleanup(func() { testutil.CleanupDatabase(t) })

		resp := testutil.ReqJSON(t, suite.Server, "GET", "/api/account/online", nil, nil)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
