package organizations

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"gitlab.com/peragrin/api/models"
	"gitlab.com/peragrin/api/service"
)

// CreateHandler saves a new organization to the database.
func (c *Config) CreateHandler(r *http.Request) *service.Response {
	account, ok := context.Get(r, "account").(models.Account)
	if !ok {
		return service.NewResponse(errAuthenticationRequired, http.StatusUnauthorized, nil)
	}

	form := models.Organization{}
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	// If there is a geocode lookup failure, then log the failure. We
	// will just let the user manually enter the coordinates.
	if err := form.SetGeo(c.LocationIQAPIKey); err != nil {
		logrus.WithFields(logrus.Fields{
			"street":  form.Street,
			"city":    form.City,
			"state":   form.State,
			"country": form.Country,
			"zip":     form.Zip,
		}).Error(errors.Wrap(err, errGeocode.Error()))
	}

	if err := form.Save(c.Client); err != nil {
		return service.NewResponse(errors.Wrap(err, errCreateOrganization.Error()), http.StatusBadRequest, nil)
	}

	if err := form.AddOperator(account.ID, c.Client); err != nil {
		return service.NewResponse(errors.Wrap(err, errAddOperator.Error()), http.StatusBadRequest, nil)
	}

	return service.NewResponse(nil, http.StatusCreated, form)
}

// UpdateHandler updates an organization.
func (c *Config) UpdateHandler(r *http.Request) *service.Response {
	form := models.Organization{}
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		return service.NewResponse(errors.Wrap(err, errUpdateOrganization.Error()), http.StatusBadRequest, nil)
	}

	id, err := strconv.Atoi(mux.Vars(r)["organizationID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errOrganizationIDRequired.Error()), http.StatusBadRequest, nil)
	}
	form.ID = id

	if err := form.Save(c.Client); err != nil {
		return service.NewResponse(errors.Wrap(err, errUpdateOrganization.Error()), http.StatusBadRequest, nil)
	}

	return service.NewResponse(nil, http.StatusOK, form)
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

// CreatePostHandler saves a new post to the database.
func (c *Config) CreatePostHandler(r *http.Request) *service.Response {
	form := models.Post{}
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		return service.NewResponse(errors.Wrap(err, errCreatePost.Error()), http.StatusBadRequest, nil)
	}

	organizationID, err := strconv.Atoi(mux.Vars(r)["organizationID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errOrganizationIDRequired.Error()), http.StatusBadRequest, nil)
	}

	form.OrganizationID = organizationID
	if err := form.Save(c.Client); err != nil {
		return service.NewResponse(errors.Wrap(err, errCreatePost.Error()), http.StatusInternalServerError, nil)
	}

	return service.NewResponse(nil, http.StatusCreated, form)
}

// ListCommunitiesHandler returns a response with all communities that are
// membered by the provided organization.
func (c *Config) ListCommunitiesHandler(r *http.Request) *service.Response {
	id, err := strconv.Atoi(mux.Vars(r)["organizationID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errOrganizationIDRequired.Error()), http.StatusBadRequest, nil)
	}

	v, err := models.ListCommunitiesByOrganizationID(id, c.Client)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, v)
}

// BulkJoinCommunityHandler creates a membership relationship between the given
// organization and communities.
func (c *Config) BulkJoinCommunityHandler(r *http.Request) *service.Response {
	id, err := strconv.Atoi(mux.Vars(r)["organizationID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errOrganizationIDRequired.Error()), http.StatusBadRequest, nil)
	}

	communities := models.Communities{}
	if err := json.NewDecoder(r.Body).Decode(&communities); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	for _, community := range communities {
		if community.AddMembership(id, false, c.Client); err != nil {
			return service.NewResponse(errors.Wrap(err, fmt.Sprintf("adding membership (o: %d, c: %d)", id, community.ID)), http.StatusBadRequest, nil)
		}
	}
	return service.NewResponse(nil, http.StatusOK, nil)
}
