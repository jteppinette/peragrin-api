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
// operated by the provided account.
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

// ListCommunitiesHandler generates a response object containing the communities that are
// conntected to the provided account.
func (c *Config) ListCommunitiesHandler(r *http.Request) *service.Response {
	id, err := strconv.Atoi(mux.Vars(r)["accountID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errAccountIDRequired.Error()), http.StatusBadRequest, nil)
	}

	communities, err := models.GetCommunitiesByAccount(id, r.URL.Query(), c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	return service.NewResponse(nil, http.StatusOK, communities)
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

// ListRedemptionsHandler returns the list of promotion redemption events for the given
// account.
func (c *Config) ListRedemptionsHandler(r *http.Request) *service.Response {
	accountID, err := strconv.Atoi(mux.Vars(r)["accountID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errAccountIDRequired.Error()), http.StatusBadRequest, nil)
	}

	redemptions, err := models.GetAccountsPromotionsByAccount(accountID, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	return service.NewResponse(nil, http.StatusOK, redemptions)
}

// ListMembershipsByCommunityHandler returns the list of memberships that an account is currently a member of.
// These memberships will be grouped into their corresponding communities.
func (c *Config) ListMembershipsByCommunityHandler(r *http.Request) *service.Response {
	accountID, err := strconv.Atoi(mux.Vars(r)["accountID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errAccountIDRequired.Error()), http.StatusBadRequest, nil)
	}

	memberships, err := models.GetMembershipsByAccount(accountID, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	if len(memberships) == 0 {
		return service.NewResponse(nil, http.StatusOK, []interface{}{})
	}

	ids := []int{}
	for _, membership := range memberships {
		ids = append(ids, membership.CommunityID)
	}

	communities, err := models.GetCommunitiesByID(ids, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	type communityWithMemberships struct {
		models.Community
		Memberships []models.Membership `json:"memberships"`
	}

	result := []communityWithMemberships{}
	for _, community := range communities {
		v := communityWithMemberships{community, []models.Membership{}}
		for _, membership := range memberships {
			if membership.CommunityID == community.ID {
				membership.CommunityID = 0
				v.Memberships = append(v.Memberships, membership)
			}
		}
		result = append(result, v)
	}

	return service.NewResponse(nil, http.StatusOK, result)
}
