package organizations

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"gitlab.com/peragrin/api/models"
	"gitlab.com/peragrin/api/service"
)

// ListHandler returns a response with all organizations.
func (c *Config) ListHandler(r *http.Request) *service.Response {
	v, err := models.ListOrganizations(c.Client)
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errListOrganizations.Error()), http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, v)
}

// GetHandler returns a response with all organizations.
func (c *Config) GetHandler(r *http.Request) *service.Response {
	id, err := strconv.Atoi(mux.Vars(r)["organizationID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errOrganizationIDRequired.Error()), http.StatusBadRequest, nil)
	}

	v, err := models.GetOrganizationByID(id, c.Client)
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errGetOrganization.Error()), http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, v)
}

// EnableHandler enabled a given organization.
func (c *Config) EnableHandler(r *http.Request) *service.Response {
	id, err := strconv.Atoi(mux.Vars(r)["organizationID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errOrganizationIDRequired.Error()), http.StatusBadRequest, nil)
	}

	v, err := models.GetOrganizationByID(id, c.Client)
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errGetOrganization.Error()), http.StatusBadRequest, nil)
	}

	v.Enabled = true
	if err := v.Save(c.Client); err != nil {
		return service.NewResponse(errors.Wrap(err, errEnableOrganization.Error()), http.StatusInternalServerError, nil)
	}

	return service.NewResponse(nil, http.StatusOK, v)
}
