package memberships

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"gitlab.com/peragrin/api/models"
	"gitlab.com/peragrin/api/service"
)

// GetHandler generates a response with the requested membership.
func (c *Config) GetHandler(r *http.Request) *service.Response {
	id, err := strconv.Atoi(mux.Vars(r)["membershipID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errMembershipIDRequired.Error()), http.StatusBadRequest, nil)
	}

	membership, err := models.GetMembershipByID(id, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	if membership == nil {
		return service.NewResponse(nil, http.StatusNotFound, nil)
	}

	return service.NewResponse(nil, http.StatusOK, membership)
}

// UpdateHandler updates the provided membership.
func (c *Config) UpdateHandler(r *http.Request) *service.Response {
	id, err := strconv.Atoi(mux.Vars(r)["membershipID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errMembershipIDRequired.Error()), http.StatusBadRequest, nil)
	}

	membership := models.Membership{}
	if err := json.NewDecoder(r.Body).Decode(&membership); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	membership.ID = id
	if err := membership.Update(c.DBClient); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	return service.NewResponse(nil, http.StatusOK, membership)
}

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

// AddAccountHandler joins a new or pre-existing account to the provided membership.
func (c *Config) AddAccountHandler(r *http.Request) *service.Response {
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
		existing.Expiration = account.Expiration
		if err := existing.AddMembership(membershipID, c.DBClient); err != nil {
			return service.NewResponse(err, http.StatusBadRequest, nil)
		}
		return service.NewResponse(nil, http.StatusOK, account)
	}

	if err := account.CreateWithMembership(membershipID, c.DBClient); err != nil {
		log.WithFields(log.Fields{
			"email": account.Email, "error": err.Error(), "membershipID": membershipID, "id": r.Header.Get("X-Request-ID"),
		}).Info(errAccountCreation.Error())
		return service.NewResponse(errAccountCreation, http.StatusBadRequest, map[string]string{"msg": errAccountCreation.Error()})
	}

	community, err := models.GetCommunityByMembershipID(membershipID, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	if err := account.SendActivationEmail(fmt.Sprintf("/map?community=%s", community.Name), c.AppDomain, c.TokenSecret, fmt.Sprintf("%s Membership", community.Name), c.MailClient); err != nil {
		log.WithFields(log.Fields{
			"email": account.Email, "error": err.Error(), "id": r.Header.Get("X-Request-ID"),
		}).Info(errAccountActivationEmail.Error())
	}

	return service.NewResponse(nil, http.StatusOK, account)
}

// UpdateAccountHandler updates an account and account membership relationship.
func (c *Config) UpdateAccountHandler(r *http.Request) *service.Response {
	account := models.Account{}
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	membershipID, err := strconv.Atoi(mux.Vars(r)["membershipID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errMembershipIDRequired.Error()), http.StatusBadRequest, nil)
	}

	accountID, err := strconv.Atoi(mux.Vars(r)["accountID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errMembershipIDRequired.Error()), http.StatusBadRequest, nil)
	}
	account.ID = accountID

	if err := account.UpdateWithMembership(membershipID, c.DBClient); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	return service.NewResponse(nil, http.StatusOK, account)
}

// DeleteHandler deletes a membership.
func (c *Config) DeleteHandler(r *http.Request) *service.Response {
	id, err := strconv.Atoi(mux.Vars(r)["membershipID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errMembershipIDRequired.Error()), http.StatusBadRequest, nil)
	}

	if err := models.DeleteMembership(id, c.DBClient); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, nil)
}

// RemoveAccountHandler removes the account - membership relationship for the provided resources.
func (c *Config) RemoveAccountHandler(r *http.Request) *service.Response {
	membershipID, err := strconv.Atoi(mux.Vars(r)["membershipID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errMembershipIDRequired.Error()), http.StatusBadRequest, nil)
	}
	accountID, err := strconv.Atoi(mux.Vars(r)["accountID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errAccountIDRequired.Error()), http.StatusBadRequest, nil)
	}

	account := models.Account{ID: accountID}
	if err := account.RemoveMembership(membershipID, c.DBClient); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	return service.NewResponse(nil, http.StatusNoContent, nil)
}
