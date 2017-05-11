package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"gitlab.com/peragrin/api/models"
	"gitlab.com/peragrin/api/service"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

var now = time.Now()

type mockClock struct{}

func (mockClock) Now() time.Time {
	return now
}

func TestLoginHandler(t *testing.T) {
	dbAccount := models.Account{Email: "jte@jte.com", Password: "jte", ID: 1}
	dbAccount.SetPassword(strings.Split(dbAccount.Email, "@")[0])

	expectedResponseToken, _ := token("secret", models.Account{Email: dbAccount.Email, ID: dbAccount.ID}, "", mockClock{})

	tests := []struct {
		bytes    []byte
		response service.Response
	}{
		{
			[]byte(`{"email": "jte@jte.com", "password": "jte"}`),
			service.Response{nil, http.StatusOK, struct {
				Token string `json:"token"`
			}{expectedResponseToken}},
		},
		{
			[]byte(`{"email": "jte@jte.com", "password": "bob"}`),
			service.Response{errInvalidCredentials, http.StatusUnauthorized, nil},
		},
		{
			[]byte(`{"email": "unknown@unknown.com", "password": "bob"}`),
			service.Response{errAccountNotFound, http.StatusUnauthorized, nil},
		},
		{
			[]byte(``),
			service.Response{errBadCredentialsFormat, http.StatusBadRequest, nil},
		},
	}

	for _, test := range tests {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		config := Init(sqlx.NewDb(db, "sqlmock"), "secret", "")
		config.Clock = mockClock{}

		creds := Credentials{}
		unmarshalErr := json.Unmarshal(test.bytes, &creds)

		if unmarshalErr == nil {
			expected := mock.ExpectQuery("^SELECT (.+) FROM Account WHERE email = (.+);").WithArgs(creds.Email)
			if creds.Email == dbAccount.Email {
				expected.WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).AddRow(dbAccount.ID, dbAccount.Email, dbAccount.Password))
			}
		}

		r, _ := http.NewRequest("", "", bytes.NewBuffer(test.bytes))
		response := config.LoginHandler(r)

		if unmarshalErr == nil {
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectation error: %s", err)
			}
		}

		validateResponse(&test.response, response, t)
	}
}

func TestRequiredMiddleware(t *testing.T) {
	dbAccount := models.Account{Email: "jte@jte.com", Password: "jte", ID: 1}
	dbAccount.SetPassword(strings.Split(dbAccount.Email, "@")[0])

	expectedResponseAccount := models.Account{Email: dbAccount.Email, ID: dbAccount.ID}

	authenticatedJWT, _ := token("secret", expectedResponseAccount, "", clock{})
	unauthenticatedJWT, _ := token("bad-secret", expectedResponseAccount, "", clock{})

	tests := []struct {
		header   http.Header
		response service.Response
	}{
		{
			http.Header{"Authorization": []string{fmt.Sprintf("Basic %s", basicAuth(dbAccount.Email, strings.Split(dbAccount.Email, "@")[0]))}},
			service.Response{nil, http.StatusOK, expectedResponseAccount},
		},
		{
			http.Header{"Authorization": []string{fmt.Sprintf("Basic %s", basicAuth(dbAccount.Email, "bad-password"))}},
			service.Response{errBasicAuth, http.StatusUnauthorized, nil},
		},
		{
			http.Header{"Authorization": []string{fmt.Sprintf("Basic %s", "bad-format")}},
			service.Response{errBadCredentialsFormat, http.StatusBadRequest, nil},
		},
		{
			http.Header{"Authorization": []string{fmt.Sprintf("Bearer %s", authenticatedJWT)}},
			service.Response{nil, http.StatusOK, expectedResponseAccount},
		},
		{
			http.Header{"Authorization": []string{fmt.Sprintf("Bearer %s", unauthenticatedJWT)}},
			service.Response{errJWTAuth, http.StatusUnauthorized, nil},
		},
		{
			http.Header{"Authorization": []string{"not-supported"}},
			service.Response{errAuthenticationStrategyNotSupported, http.StatusUnauthorized, nil},
		},
		{
			http.Header{},
			service.Response{errAuthenticationRequired, http.StatusUnauthorized, nil},
		},
	}

	for _, test := range tests {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		config := Init(sqlx.NewDb(db, "sqlmock"), "secret", "")

		var expected *sqlmock.ExpectedQuery
		if email, _, ok := parseBasicAuth(test.header.Get("Authorization")); ok {
			expected = mock.ExpectQuery("^SELECT (.+) FROM Account WHERE email = (.+);").WithArgs(email)
			expected.WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).AddRow(dbAccount.ID, dbAccount.Email, dbAccount.Password))
		}

		response := config.RequiredMiddleware(config.AccountHandler)(&http.Request{Header: test.header})

		if err := mock.ExpectationsWereMet(); expected != nil && err != nil {
			t.Errorf("unmet expectation error: %s", err)
		}

		validateResponse(&test.response, response, t)
	}
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}

func getError(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func validateResponse(expected *service.Response, actual *service.Response, t *testing.T) {
	if actual.Code != expected.Code {
		t.Errorf("expected response code to be %d, got %d", expected.Code, actual.Code)
	}
	if !reflect.DeepEqual(actual.Data, expected.Data) {
		t.Errorf("expected response data to be %v, got %v", expected.Data, actual.Data)
	}
	if !strings.Contains(getError(actual.Error), getError(expected.Error)) {
		t.Errorf("expected response error to contain '%s', got '%s'", getError(expected.Error), getError(actual.Error))
	}
}
