package communities

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"gitlab.com/peragrin/api/models"
	"gitlab.com/peragrin/api/service"
)

// ListHandler returns a response with all communities.
func (c *Config) ListHandler(r *http.Request) *service.Response {
	communities, err := models.GetCommunities(c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	return service.NewResponse(nil, http.StatusOK, communities)
}

// GetHandler returns a response with the requested community.
func (c *Config) GetHandler(r *http.Request) *service.Response {
	id, err := strconv.Atoi(mux.Vars(r)["communityID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errCommunityIDRequired.Error()), http.StatusBadRequest, nil)
	}

	community, err := models.GetCommunityByID(id, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	return service.NewResponse(nil, http.StatusOK, community)
}

// ListOrganizationsHandler returns a response with all organizations
// in a given community.
func (c *Config) ListOrganizationsHandler(r *http.Request) *service.Response {
	communityID, err := strconv.Atoi(mux.Vars(r)["communityID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errCommunityIDRequired.Error()), http.StatusBadRequest, nil)
	}

	organizations, err := models.GetOrganizationsByCommunity(communityID, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	if err := organizations.SetPresignedLogoLinks(c.StoreClient); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	return service.NewResponse(nil, http.StatusOK, organizations)
}

// CreateOrganizationHandler creates a new organization that is automatically
// joined with the creating community.
func (c *Config) CreateOrganizationHandler(r *http.Request) *service.Response {
	communityID, err := strconv.Atoi(mux.Vars(r)["communityID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errCommunityIDRequired.Error()), http.StatusBadRequest, nil)
	}

	organization := models.Organization{}
	if err := json.NewDecoder(r.Body).Decode(&organization); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	if err := organization.CreateWithCommunity(communityID, c.DBClient); err != nil {
		return service.NewResponse(errors.Wrap(err, errCreateOrganization.Error()), http.StatusBadRequest, nil)
	}

	return service.NewResponse(nil, http.StatusCreated, organization)
}

// ListMembershipsHandler returns a response with all memberships
// in a given community.
func (c *Config) ListMembershipsHandler(r *http.Request) *service.Response {
	communityID, err := strconv.Atoi(mux.Vars(r)["communityID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errCommunityIDRequired.Error()), http.StatusBadRequest, nil)
	}

	memberships, err := models.GetMembershipsByCommunity(communityID, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	return service.NewResponse(nil, http.StatusOK, memberships)
}

// CreateMembershipHandler saves a new membership to the database.
func (c *Config) CreateMembershipHandler(r *http.Request) *service.Response {
	membership := models.Membership{}
	if err := json.NewDecoder(r.Body).Decode(&membership); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	communityID, err := strconv.Atoi(mux.Vars(r)["communityID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errCommunityIDRequired.Error()), http.StatusBadRequest, nil)
	}

	if err := membership.Create(communityID, c.DBClient); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusCreated, membership)
}

// ListGeoJSONOverlaysHandler returns a response with all geo JSON overlays
// in a given community.
func (c *Config) ListGeoJSONOverlaysHandler(r *http.Request) *service.Response {
	communityID, err := strconv.Atoi(mux.Vars(r)["communityID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errCommunityIDRequired.Error()), http.StatusBadRequest, nil)
	}

	geoJSONOverlays, err := models.GetGeoJSONOverlaysByCommunity(communityID, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, geoJSONOverlays)
}

// ListPostsHandler returns a response with all posts
// in a given community.
func (c *Config) ListPostsHandler(r *http.Request) *service.Response {
	communityID, err := strconv.Atoi(mux.Vars(r)["communityID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errCommunityIDRequired.Error()), http.StatusBadRequest, nil)
	}

	posts, err := models.GetPostsByCommunity(communityID, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, posts)
}

// BulkAddAccountsHandler creates and invites multiple accounts to a community as business operators
// in a single atomic action.
func (c *Config) BulkAddAccountsHandler(r *http.Request) *service.Response {
	id, err := strconv.Atoi(mux.Vars(r)["communityID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errCommunityIDRequired.Error()), http.StatusBadRequest, nil)
	}

	community, err := models.GetCommunityByID(id, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	accounts := models.Accounts{}
	if err := json.NewDecoder(r.Body).Decode(&accounts); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	emails := []string{}
	for _, account := range accounts {
		emails = append(emails, account.Email)
	}

	existing, err := models.GetAccountsByEmails(emails, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	needs := models.Accounts{}
	for _, account := range accounts {
		var exists bool
		for _, lookup := range existing {
			if account.Email == lookup.Email {
				exists = true
				break
			}
		}
		if !exists {
			needs = append(needs, account)
		}
	}
	if err := needs.Create(c.DBClient); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.WithFields(log.Fields{"error": err, "id": r.Header.Get("X-Request-ID")}).Info(errAccountActivationEmail.Error())
			}
		}()
		for _, need := range needs {
			if err := need.SendActivationEmail("/setup/business-leader", c.AppDomain, c.TokenSecret, fmt.Sprintf("%s Business Operator", community.Name), c.MailClient); err != nil {
				log.WithFields(log.Fields{
					"email": need.Email, "error": err.Error(), "id": r.Header.Get("X-Request-ID"),
				}).Info(errAccountActivationEmail.Error())
			}
		}
	}()

	return service.NewResponse(nil, http.StatusOK, nil)
}

// UpdateHandler updates an community.
func (c *Config) UpdateHandler(r *http.Request) *service.Response {
	id, err := strconv.Atoi(mux.Vars(r)["communityID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errCommunityIDRequired.Error()), http.StatusBadRequest, nil)
	}

	community := models.Community{}
	if err := json.NewDecoder(r.Body).Decode(&community); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	community.ID = id

	if err := community.Update(c.DBClient); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, community)
}

// CreateHandler creates a new community. This requires that the requesting account is
// a super user.
func (c *Config) CreateHandler(r *http.Request) *service.Response {
	account, ok := context.Get(r, "account").(models.Account)
	if !ok {
		return service.NewResponse(errAuthenticationRequired, http.StatusUnauthorized, nil)
	}
	if !account.IsSuper {
		return service.NewResponse(errSuperUserRequired, http.StatusForbidden, nil)
	}

	community := models.Community{}
	if err := json.NewDecoder(r.Body).Decode(&community); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	if err := community.Create(c.DBClient); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	return service.NewResponse(nil, http.StatusCreated, community)
}
