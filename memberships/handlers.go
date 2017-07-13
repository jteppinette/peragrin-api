package memberships

import (
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

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

	accounts, err := models.GetAccountsByMembership(membershipID, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, accounts)
}

// CreateAccountHandler creates a new account and connects it to the
// provided membership.
func (c *Config) CreateAccountHandler(r *http.Request) *service.Response {
	account := models.Account{}
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	membershipID, err := strconv.Atoi(mux.Vars(r)["membershipID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errMembershipIDRequired.Error()), http.StatusBadRequest, nil)
	}

	// If the account already exists, then simply add the membership to it.
	if existing, err := models.GetAccountByEmail(account.Email, c.DBClient); err != nil {
		return service.NewResponse(nil, http.StatusBadRequest, nil)
	} else if existing != nil {
		if err := existing.AddMembership(membershipID, c.DBClient); err != nil {
			return service.NewResponse(err, http.StatusBadRequest, nil)
		}
		return service.NewResponse(nil, http.StatusOK, account)
	}

	account.Password = ""
	if err := account.CreateWithMembership(membershipID, c.DBClient); err != nil {
		log.WithFields(log.Fields{
			"email": account.Email, "error": err.Error(), "membershipID": membershipID, "id": r.Header.Get("X-Request-ID"),
		}).Info(errAccountCreationFailed.Error())
		return service.NewResponse(errAccountCreationFailed, http.StatusBadRequest, map[string]string{"msg": errAccountCreationFailed.Error()})
	}

	if err := account.SendAccountActivationEmail(c.AppDomain, c.TokenSecret, c.Clock, c.MailClient); err != nil {
		log.WithFields(log.Fields{
			"email": account.Email, "error": err.Error(), "id": r.Header.Get("X-Request-ID"),
		}).Info(errAccountActivationEmail.Error())
	}

	return service.NewResponse(nil, http.StatusOK, account)
}
