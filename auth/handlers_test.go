package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"gitlab.com/peragrin/api/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestLoginHandler(t *testing.T) {
	tts := []lhtt{
		lhtt{Credentials{"jteppinette", "jteppinette"}, models.User{ID: 1, Username: "jteppinette", OrganizationID: 1}, http.StatusOK},
		lhtt{Credentials{"jteppinette", "bad"}, models.User{ID: 1, Username: "jteppinette", OrganizationID: 1}, http.StatusUnauthorized},
	}
	for _, tt := range tts {
		tt.test(t)
	}
}

type lhtt struct {
	creds Credentials
	user  models.User
	code  int
}

func (v lhtt) test(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	config := Init(sqlx.NewDb(db, "sqlmock"), "secret")

	v.user.SetPassword(v.user.Username)
	columns := []string{"id", "username", "password", "organizationid"}
	mock.ExpectQuery("^SELECT (.+) FROM users WHERE username = (.+);").
		WithArgs(v.creds.Username).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(v.user.ID, v.user.Username, v.user.Password, v.user.OrganizationID))

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(v.creds)
	r, _ := http.NewRequest("POST", "", &b)
	w := httptest.NewRecorder()
	config.LoginHandler(w, r)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
	}

	if w.Code != v.code {
		t.Errorf("expected code to be %d, got %d:%v", v.code, w.Code, w.Body)
	}

	// If we are testing a failure case, then return early.
	if v.code != http.StatusOK {
		return
	}

	var au AuthUser
	json.NewDecoder(w.Body).Decode(&au)
	v.user.Password = ""
	if !reflect.DeepEqual(au.User, v.user) {
		t.Errorf("expected user to be %v, got %v", v.user, au.User)
	}
	if au.Token == "" {
		t.Errorf("expected token to not be '', got %s", au.Token)
	}
}
