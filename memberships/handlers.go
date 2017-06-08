package memberships

import (
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"gitlab.com/peragrin/api/auth"
	"gitlab.com/peragrin/api/models"
	"gitlab.com/peragrin/api/service"
)

// ListAccountsHandler returns a response with all accounts that are
// have the provided membership.
func (c *Config) ListAccountsHandler(r *http.Request) *service.Response {
	membershipID, err := strconv.Atoi(mux.Vars(r)["membershipID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errMembershipIDRequired.Error()), http.StatusBadRequest, nil)
	}

	accounts, err := models.GetAccountsByMembership(membershipID, c.Client)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, accounts)
}

// CreateAccountHandler creates a new account and connects it to the
// provided membership.
func (c *Config) CreateAccountHandler(r *http.Request) *service.Response {
	creds := auth.Credentials{}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	membershipID, err := strconv.Atoi(mux.Vars(r)["membershipID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errMembershipIDRequired.Error()), http.StatusBadRequest, nil)
	}

	account := &models.Account{Email: creds.Email}
	account.SetPassword(creds.Password)
	if err := account.CreateWithMembership(membershipID, c.Client); err != nil {
		log.WithFields(log.Fields{
			"email": creds.Email, "error": err.Error(), "membershipID": membershipID, "id": r.Header.Get("X-Request-ID"),
		}).Info(errAccountCreationFailed.Error())
		return service.NewResponse(errAccountCreationFailed, http.StatusBadRequest, map[string]string{"msg": errAccountCreationFailed.Error()})
	}

	return service.NewResponse(nil, http.StatusOK, account)
}
