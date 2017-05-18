package communities

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"gitlab.com/peragrin/api/models"
	"gitlab.com/peragrin/api/service"
)

// ListHandler returns a response with all communities.
func (c *Config) ListHandler(r *http.Request) *service.Response {
	communities, err := models.GetCommunities(c.Client)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, communities)
}

// ListOrganizationsHandler returns a response with all organizations
// in a given community.
func (c *Config) ListOrganizationsHandler(r *http.Request) *service.Response {
	communityID, err := strconv.Atoi(mux.Vars(r)["communityID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errCommunityIDRequired.Error()), http.StatusBadRequest, nil)
	}

	organizations, err := models.GetOrganizationsByCommunity(communityID, c.Client)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, organizations)
}

// ListPostsHandler returns a response with all posts
// in a given community.
func (c *Config) ListPostsHandler(r *http.Request) *service.Response {
	communityID, err := strconv.Atoi(mux.Vars(r)["communityID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errCommunityIDRequired.Error()), http.StatusBadRequest, nil)
	}

	posts, err := models.GetPostsByCommunity(communityID, c.Client)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, posts)
}
