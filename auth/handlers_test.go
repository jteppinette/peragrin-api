package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.com/peragrin/api/db"
	"gitlab.com/peragrin/api/fixture"
)

var config *Config

func init() {
	client, err := db.Client("0.0.0.0", "db", "secret", "db")
	if err != nil {
		panic(err)
	}
	db.Migrate(client)
	fixture.Initialize(client)

	config = Init(client, "secret")
}

func TestLoginHandler(t *testing.T) {
	creds := Credentials{"jteppinette", "jteppinette"}

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(creds)
	r, _ := http.NewRequest("POST", "", &b)
	w := httptest.NewRecorder()
	config.LoginHandler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected code to be 200, got %d", w.Code)
	}

	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)
	if result["username"] != creds.Username {
		t.Errorf("expected username to be %s, got %v", creds.Username, result["username"])
	}
	if result["token"] == nil {
		t.Error("expected token to be non-nil, got nil")
	}
	if result["organizationID"] == nil {
		t.Error("expected organizationID to be non-nil, got nil")
	}
	if result["id"] == nil {
		t.Error("expected id to be non-nil, got nil")
	}
}
