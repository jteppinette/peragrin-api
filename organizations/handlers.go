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

// UpdateHandler updates an organization.
func (c *Config) UpdateHandler(r *http.Request) *service.Response {
	organization := models.Organization{}
	if err := json.NewDecoder(r.Body).Decode(&organization); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	id, err := strconv.Atoi(mux.Vars(r)["organizationID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errOrganizationIDRequired.Error()), http.StatusBadRequest, nil)
	}
	organization.ID = id

	if err := organization.Update(c.DBClient); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, organization)
}

// GetHandler generates a response with the requested organization.
func (c *Config) GetHandler(r *http.Request) *service.Response {
	id, err := strconv.Atoi(mux.Vars(r)["organizationID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errOrganizationIDRequired.Error()), http.StatusBadRequest, nil)
	}

	organization, err := models.GetOrganizationByID(id, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	if err := organization.SetPresignedLogoLink(c.StoreClient); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	return service.NewResponse(nil, http.StatusOK, organization)
}

// CreatePostHandler saves a new post to the database.
func (c *Config) CreatePostHandler(r *http.Request) *service.Response {
	post := models.Post{}
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	organizationID, err := strconv.Atoi(mux.Vars(r)["organizationID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errOrganizationIDRequired.Error()), http.StatusBadRequest, nil)
	}

	post.OrganizationID = organizationID
	if err := post.Save(c.DBClient); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusCreated, post)
}

// CreatePromotionHandler saves a new promotion to the database.
func (c *Config) CreatePromotionHandler(r *http.Request) *service.Response {
	promotion := models.Promotion{}
	if err := json.NewDecoder(r.Body).Decode(&promotion); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	organizationID, err := strconv.Atoi(mux.Vars(r)["organizationID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errOrganizationIDRequired.Error()), http.StatusBadRequest, nil)
	}

	promotion.OrganizationID = organizationID
	if err := promotion.Save(c.DBClient); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusCreated, promotion)
}

// ListCommunitiesHandler returns a response with all communities that are
// membered by the provided organization.
func (c *Config) ListCommunitiesHandler(r *http.Request) *service.Response {
	organizationID, err := strconv.Atoi(mux.Vars(r)["organizationID"])
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	communities, err := models.GetCommunitiesByOrganization(organizationID, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, communities)
}

// JoinCommunityHandler creates a relationship between the given
// organization and community.
func (c *Config) JoinCommunityHandler(r *http.Request) *service.Response {
	organizationID, err := strconv.Atoi(mux.Vars(r)["organizationID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errOrganizationIDRequired.Error()), http.StatusBadRequest, nil)
	}

	communityID, err := strconv.Atoi(mux.Vars(r)["communityID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errCommunityIDRequired.Error()), http.StatusBadRequest, nil)
	}

	co := models.CommunityOrganization{CommunityID: communityID, OrganizationID: organizationID}
	if co.Create(c.DBClient); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, nil)
}

// CreateCommunityHandler creates a new community and creates an
// administrative relationship between the requesting
// organization and community.
func (c *Config) CreateCommunityHandler(r *http.Request) *service.Response {
	organizationID, err := strconv.Atoi(mux.Vars(r)["organizationID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errOrganizationIDRequired.Error()), http.StatusBadRequest, nil)
	}

	community := models.Community{}
	if err := json.NewDecoder(r.Body).Decode(&community); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	if err := community.Create(organizationID, c.DBClient); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusCreated, community)
}

// SetHoursHandler sets the given hours for the requested organization.
func (c *Config) SetHoursHandler(r *http.Request) *service.Response {
	organizationID, err := strconv.Atoi(mux.Vars(r)["organizationID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errOrganizationIDRequired.Error()), http.StatusBadRequest, nil)
	}

	hours := models.Hours{}
	if err := json.NewDecoder(r.Body).Decode(&hours); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	if err := hours.Set(organizationID, c.DBClient); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, nil)
}

// ListHoursHandler generates a response with the operational hours for
// the requested organization.
func (c *Config) ListHoursHandler(r *http.Request) *service.Response {
	organizationID, err := strconv.Atoi(mux.Vars(r)["organizationID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errOrganizationIDRequired.Error()), http.StatusBadRequest, nil)
	}

	hours, err := models.GetHoursByOrganization(organizationID, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	return service.NewResponse(nil, http.StatusOK, hours)
}

// ListPromotionsHandler generates a response with the promotions for
// the requested organization.
func (c *Config) ListPromotionsHandler(r *http.Request) *service.Response {
	organizationID, err := strconv.Atoi(mux.Vars(r)["organizationID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errOrganizationIDRequired.Error()), http.StatusBadRequest, nil)
	}

	promotions, err := models.GetPromotionsByOrganization(organizationID, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	return service.NewResponse(nil, http.StatusOK, promotions)
}
