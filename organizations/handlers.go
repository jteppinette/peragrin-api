package organizations

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"gitlab.com/peragrin/api/models"
	"gitlab.com/peragrin/api/service"
)

// CreateHandler saves a new organization to the database.
func (c *Config) CreateHandler(r *http.Request) *service.Response {
	form := models.Organization{}
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		return service.NewResponse(errors.Wrap(err, errCreateOrganization.Error()), http.StatusBadRequest, nil)
	}

	if err := form.SetGeo(form.Address, c.LocationIQAPIKey); err != nil {
		return service.NewResponse(errors.Wrap(errors.Wrap(err, errGeocode.Error()), errCreateOrganization.Error()), http.StatusBadRequest, nil)
	}

	if err := form.Save(c.Client); err != nil {
		return service.NewResponse(errors.Wrap(err, errCreateOrganization.Error()), http.StatusBadRequest, nil)
	}

	return service.NewResponse(nil, http.StatusCreated, form)
}

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
