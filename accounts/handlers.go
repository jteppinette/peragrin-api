package accounts

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"gitlab.com/peragrin/api/models"
	"gitlab.com/peragrin/api/service"
)

// UpdateAccountHandler updates an account.
func (c *Config) UpdateAccountHandler(r *http.Request) *service.Response {
	id, err := strconv.Atoi(mux.Vars(r)["accountID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errAccountIDRequired.Error()), http.StatusBadRequest, nil)
	}

	account := models.Account{}
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	account.ID = id

	if err := account.Save(c.DBClient); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, account)
}

// ForgotPasswordHandler generates a token that can be used to reset the password
// of the account with the provided email address.
func (c *Config) ForgotPasswordHandler(r *http.Request) *service.Response {
	id, err := strconv.Atoi(mux.Vars(r)["accountID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errAccountIDRequired.Error()), http.StatusBadRequest, nil)
	}

	account, err := models.GetAccountByID(id, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	if account == nil {
		return service.NewResponse(errAccountNotFound, http.StatusNotFound, nil)
	}

	if err := account.SendResetPasswordEmail(c.AppDomain, c.TokenSecret, c.MailClient); err != nil {
		return service.NewResponse(err, http.StatusInternalServerError, nil)
	}
	return service.NewResponse(nil, http.StatusOK, nil)
}

// ListOrganizationsHandler generates a response object containing the organizations that are
// operated by the currently authenticated account.
func (c *Config) ListOrganizationsHandler(r *http.Request) *service.Response {
	id, err := strconv.Atoi(mux.Vars(r)["accountID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errAccountIDRequired.Error()), http.StatusBadRequest, nil)
	}

	organizations, err := models.GetOrganizationsByAccount(id, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	if err := organizations.SetPresignedLogoLinks(c.StoreClient); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, organizations)
}

// CreateOrganizationHandler saves a new organization to the database.
func (c *Config) CreateOrganizationHandler(r *http.Request) *service.Response {
	id, err := strconv.Atoi(mux.Vars(r)["accountID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errAccountIDRequired.Error()), http.StatusBadRequest, nil)
	}

	organization := models.Organization{}
	if err := json.NewDecoder(r.Body).Decode(&organization); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	if err := organization.CreateWithAccount(id, c.DBClient); err != nil {
		return service.NewResponse(errors.Wrap(err, errCreateOrganization.Error()), http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusCreated, organization)
}

// ListHandler returns a response with all accounts.
func (c *Config) ListHandler(r *http.Request) *service.Response {
	values := r.URL.Query()

	email := values.Get("email")
	if email == "" {
		return service.NewResponse(nil, http.StatusBadRequest, nil)
	}
	account, err := models.GetAccountByEmail(email, c.DBClient)
	if err != nil {
		return service.NewResponse(nil, http.StatusBadRequest, nil)
	}

	type response struct {
		Results models.Accounts `json:"results"`
		Total   int             `json:"total"`
	}

	if account == nil {
		return service.NewResponse(nil, http.StatusOK, response{models.Accounts{}, 0})
	}
	return service.NewResponse(nil, http.StatusOK, response{models.Accounts{*account}, 1})
}

// ListPromotionRedemptionsHandler returns the list of promotion redemption events for the given
// account and promotion.
func (c *Config) ListPromotionRedemptionsHandler(r *http.Request) *service.Response {
	accountID, err := strconv.Atoi(mux.Vars(r)["accountID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errAccountIDRequired.Error()), http.StatusBadRequest, nil)
	}

	promotionID, err := strconv.Atoi(mux.Vars(r)["promotionID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errPromotionIDRequired.Error()), http.StatusBadRequest, nil)
	}

	events, err := models.GetAccountsPromotionsByID(accountID, promotionID, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	return service.NewResponse(nil, http.StatusOK, events)
}
