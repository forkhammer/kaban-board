package testutil

import (
	"main/internal/testutil/factories"
	"net/http/httptest"
	"testing"
)

func GetAdminToken(t *testing.T, server *httptest.Server) string {
	t.Helper()

	factory := factories.NewAccountFactory()
	account, err := factory.Create(map[string]any{
		"password":  "adminpass123",
		"is_active": true,
		"role":      "admin",
	})
	if err != nil {
		t.Fatalf("failed to create admin account: %v", err)
	}

	loginBody := map[string]string{
		"username": account.Username,
		"password": "adminpass123",
	}
	resp := ReqJSON(t, server, "POST", "/api/account/login", loginBody, nil)

	var loginResult map[string]any
	ParseBody(t, resp, &loginResult)

	token, ok := loginResult["token"].(string)
	if !ok || token == "" {
		t.Fatalf("failed to get admin token: %v", loginResult)
	}

	return token
}

func GetViewerToken(t *testing.T, server *httptest.Server) string {
	t.Helper()

	factory := factories.NewAccountFactory()
	account, err := factory.Create(map[string]any{
		"password":  "viewerpass123",
		"is_active": true,
		"role":      "viewer",
	})
	if err != nil {
		t.Fatalf("failed to create viewer account: %v", err)
	}

	loginBody := map[string]string{
		"username": account.Username,
		"password": "viewerpass123",
	}
	resp := ReqJSON(t, server, "POST", "/api/account/login", loginBody, nil)

	var loginResult map[string]any
	ParseBody(t, resp, &loginResult)

	token, ok := loginResult["token"].(string)
	if !ok || token == "" {
		t.Fatalf("failed to get viewer token: %v", loginResult)
	}

	return token
}
