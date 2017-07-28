package communities

import (
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
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

	// If there was no assigned geo coordinates, then run the geocode lookup.
	if organization.Lon == 0 || organization.Lat == 0 {

		// If there is a geocode lookup failure, then log the failure. We
		// will just let the user manually enter the coordinates.
		if err := organization.SetGeo(c.LocationIQAPIKey); err != nil {
			log.WithFields(log.Fields{
				"street": organization.Street, "city": organization.City, "state": organization.State, "country": organization.Country, "zip": organization.Zip,
				"error": err.Error(),
				"id":    r.Header.Get("X-Request-ID"),
			}).Info(errGeocode.Error())
		}
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
