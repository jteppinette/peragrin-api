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
	v, err := models.ListCommunities(c.Client)
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errListCommunities.Error()), http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, v)
}

// ListOrganizationsHandler returns a response with all organizations
// in a given community..
func (c *Config) ListOrganizationsHandler(r *http.Request) *service.Response {
	id, err := strconv.Atoi(mux.Vars(r)["communityID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errCommunityIDRequired.Error()), http.StatusBadRequest, nil)
	}

	v, err := models.ListOrganizationsByCommunityID(id, c.Client)
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errListOrganizations.Error()), http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, v)
}
