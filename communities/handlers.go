package communities

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"gitlab.com/peragrin/api/models"
	"gitlab.com/peragrin/api/service"
)

// ListHandler returns a response with all communities.
func (c *Config) ListHandler(r *http.Request) *service.Response {
	debug(r, "initialize: ListHandler")

	communities, err := models.GetCommunities(c.DBClient)
	if err != nil {
		debug(r, "error: ListHandler")
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	debug(r, "return: ListHandler")
	return service.NewResponse(nil, http.StatusOK, communities)
}

// ListOrganizationsHandler returns a response with all organizations
// in a given community.
func (c *Config) ListOrganizationsHandler(r *http.Request) *service.Response {
	debug(r, "initialize: ListHandler")

	communityID, err := strconv.Atoi(mux.Vars(r)["communityID"])
	if err != nil {
		debug(r, fmt.Sprintf("error: ListOrganizationsHandler: %s", err.Error()))
		return service.NewResponse(errors.Wrap(err, errCommunityIDRequired.Error()), http.StatusBadRequest, nil)
	}

	organizations, err := models.GetOrganizationsByCommunity(communityID, c.DBClient)
	if err != nil {
		debug(r, fmt.Sprintf("error: ListOrganizationsHandler: %s", err.Error()))
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	if err := organizations.SetPresignedLogoLinks(c.StoreClient); err != nil {
		debug(r, fmt.Sprintf("error: ListOrganizationsHandler: %s", err.Error()))
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	if err := organizations.SetPresignedIconLinks(c.StoreClient); err != nil {
		debug(r, fmt.Sprintf("error: ListOrganizationsHandler: %s", err.Error()))
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	debug(r, "return: ListOrganizationsHandler")
	return service.NewResponse(nil, http.StatusOK, organizations)
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

	if err := membership.Save(communityID, c.DBClient); err != nil {
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
