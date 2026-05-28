package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/config"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goioc/di"
)

func CleanupDatabase(t *testing.T) {
	t.Helper()

	conn := di.GetInstance("connection").(interfaces.ConnectionInterface)
	cfg := di.GetInstance("config").(*config.Config)
	db := conn.GetEngine()

	tables := []string{
		"issue_binding_histories",
		"issue_bindings",
		"label_histories",
		"sprint_user_settings",
		"assignees",
		"issue_labels",
		"issues",
		"epics",
		"releases",
		"sprints",
		"columns",
		"projects",
		"team_groups",
		"user_groups",
		"gitlab_tokens",
		"teams",
		"users",
		"groups",
		"labels",
		"accounts",
		"kv_elements",
	}

	quoteChar := "`"
	if cfg.DbType == "postgresql" {
		quoteChar = `"`
	}

	for _, table := range tables {
		sql := fmt.Sprintf("DELETE FROM %s%s%s", quoteChar, table, quoteChar)
		if err := db.Exec(sql).Error; err != nil {
			t.Logf("Warning: failed to cleanup table %s: %v", table, err)
		}
	}
}

func DoJSON(t *testing.T, server *httptest.Server, method, path string, body any) *http.Response {
	t.Helper()

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, server.URL+path, reqBody)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}

	return resp
}

func DoAuthJSON(t *testing.T, server *httptest.Server, method, path string, body any, token string) *http.Response {
	t.Helper()

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, server.URL+path, reqBody)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}

	return resp
}

func ParseBody(t *testing.T, resp *http.Response, result any) {
	t.Helper()

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if err := json.Unmarshal(bodyBytes, result); err != nil {
		t.Fatalf("failed to unmarshal response body: %v, body: %s", err, string(bodyBytes))
	}
}

func ActivateUser(t *testing.T, username string) {
	t.Helper()

	conn := di.GetInstance("connection").(interfaces.ConnectionInterface)
	if err := conn.GetEngine().Model(&models.Account{}).Where("username = ?", username).Update("is_active", true).Error; err != nil {
		t.Fatalf("failed to activate user: %v", err)
	}
}
