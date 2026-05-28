package integration

import (
	"main/internal/testutil"
	"net/http"
	"testing"
)

func TestRegister(t *testing.T) {
	t.Run("valid registration", func(t *testing.T) {
		t.Cleanup(func() { testutil.CleanupDatabase(t) })
		
		body := map[string]string{
			"username": "testuser",
			"password": "testpass123",
		}

		resp := testutil.DoJSON(t, suite.Server, "POST", "/api/account/register", body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		if result["username"] != "testuser" {
			t.Errorf("expected username 'testuser', got %v", result["username"])
		}
	})

	t.Run("missing username", func(t *testing.T) {
		t.Cleanup(func() { testutil.CleanupDatabase(t) })
		
		body := map[string]string{
			"password": "testpass123",
		}

		resp := testutil.DoJSON(t, suite.Server, "POST", "/api/account/register", body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("missing password", func(t *testing.T) {
		t.Cleanup(func() { testutil.CleanupDatabase(t) })
		
		body := map[string]string{
			"username": "testuser2",
		}

		resp := testutil.DoJSON(t, suite.Server, "POST", "/api/account/register", body)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", resp.StatusCode)
		}
	})
}

func TestLogin(t *testing.T) {
	t.Run("valid login", func(t *testing.T) {
		t.Cleanup(func() { testutil.CleanupDatabase(t) })
		
		registerBody := map[string]string{
			"username": "loginuser",
			"password": "loginpass123",
		}
		resp := testutil.DoJSON(t, suite.Server, "POST", "/api/account/register", registerBody)
		resp.Body.Close()

		testutil.ActivateUser(t, "loginuser")

		loginBody := map[string]string{
			"username": "loginuser",
			"password": "loginpass123",
		}
		resp = testutil.DoJSON(t, suite.Server, "POST", "/api/account/login", loginBody)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		if result["token"] == nil || result["token"] == "" {
			t.Error("expected token in response")
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		t.Cleanup(func() { testutil.CleanupDatabase(t) })
		
		registerBody := map[string]string{
			"username": "wrongpassuser",
			"password": "correctpass",
		}
		resp := testutil.DoJSON(t, suite.Server, "POST", "/api/account/register", registerBody)
		resp.Body.Close()

		loginBody := map[string]string{
			"username": "wrongpassuser",
			"password": "wrongpass",
		}
		resp = testutil.DoJSON(t, suite.Server, "POST", "/api/account/login", loginBody)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			t.Error("expected non-200 status for wrong password")
		}
	})

	t.Run("nonexistent user", func(t *testing.T) {
		t.Cleanup(func() { testutil.CleanupDatabase(t) })
		
		loginBody := map[string]string{
			"username": "nonexistent",
			"password": "somepass",
		}
		resp := testutil.DoJSON(t, suite.Server, "POST", "/api/account/login", loginBody)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			t.Error("expected non-200 status for nonexistent user")
		}
	})
}

func TestGetActiveUser(t *testing.T) {
	t.Run("with valid token", func(t *testing.T) {
		t.Cleanup(func() { testutil.CleanupDatabase(t) })
		
		registerBody := map[string]string{
			"username": "activeuser",
			"password": "activepass123",
		}
		resp := testutil.DoJSON(t, suite.Server, "POST", "/api/account/register", registerBody)
		resp.Body.Close()

		testutil.ActivateUser(t, "activeuser")

		loginBody := map[string]string{
			"username": "activeuser",
			"password": "activepass123",
		}
		resp = testutil.DoJSON(t, suite.Server, "POST", "/api/account/login", loginBody)
		var loginResult map[string]any
		testutil.ParseBody(t, resp, &loginResult)
		resp.Body.Close()

		token := loginResult["token"].(string)

		resp = testutil.DoAuthJSON(t, suite.Server, "GET", "/api/account/user", nil, token)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		user := result["user"].(map[string]any)
		if user["username"] != "activeuser" {
			t.Errorf("expected username 'activeuser', got %v", user["username"])
		}
	})

	t.Run("without token", func(t *testing.T) {
		t.Cleanup(func() { testutil.CleanupDatabase(t) })
		
		resp := testutil.DoJSON(t, suite.Server, "GET", "/api/account/user", nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]any
		testutil.ParseBody(t, resp, &result)

		if result["user"] != nil {
			t.Errorf("expected user to be null, got %v", result["user"])
		}
	})
}

func TestProtectedRoute(t *testing.T) {
	t.Run("without auth returns 401", func(t *testing.T) {
		t.Cleanup(func() { testutil.CleanupDatabase(t) })
		
		resp := testutil.DoJSON(t, suite.Server, "GET", "/api/account/online", nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", resp.StatusCode)
		}
	})
}
