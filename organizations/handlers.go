package organizations

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

	// If this organization does not have lon/lat, then query the location service.
	if organization.Lon == 0 || organization.Lat == 0 {
		if err := organization.SetGeo(c.LocationIQAPIKey); err != nil {
			log.WithFields(log.Fields{
				"street": organization.Street, "city": organization.City, "state": organization.State, "country": organization.Country, "zip": organization.Zip,
				"error": err.Error(),
				"id":    r.Header.Get("X-Request-ID"),
			}).Info(errGeocode.Error())
		}
	}

	if err := organization.Update(c.DBClient); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, organization)
}

// UploadLogoHandler uploads a new logo to the store and sets the
// organization's logo field.
func (c *Config) UploadLogoHandler(r *http.Request) *service.Response {
	id, err := strconv.Atoi(mux.Vars(r)["organizationID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errOrganizationIDRequired.Error()), http.StatusBadRequest, nil)
	}

	file, header, err := r.FormFile("logo")
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	organization, err := models.GetOrganizationByID(id, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusInternalServerError, nil)
	}

	organization.Hours, err = models.GetHoursByOrganization(organization.ID, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	organization.Logo = header.Filename
	if err := organization.UploadLogo(file, c.StoreClient); err != nil {
		log.WithFields(log.Fields{
			"file":   file,
			"header": header,
			"error":  err.Error(),
			"id":     r.Header.Get("X-Request-ID"),
		}).Info(errUploadLogo.Error())
		return service.NewResponse(errUploadLogo, http.StatusBadRequest, nil)
	}

	if err := organization.Update(c.DBClient); err != nil {
		log.WithFields(log.Fields{
			"logo":           organization.Logo,
			"organizationID": organization.ID,
			"error":          err.Error(),
			"id":             r.Header.Get("X-Request-ID"),
		}).Info(errUpdateOrganization.Error())
		return service.NewResponse(errUpdateOrganization, http.StatusBadRequest, nil)
	}

	if err := organization.SetPresignedLogoLink(c.StoreClient); err != nil {
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

	organization.Hours, err = models.GetHoursByOrganization(organization.ID, c.DBClient)
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

// ListAccountsHandler returns a response with all accounts that are
// operating the given organization.
func (c *Config) ListAccountsHandler(r *http.Request) *service.Response {
	organizationID, err := strconv.Atoi(mux.Vars(r)["organizationID"])
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	accounts, err := models.GetAccountsByOrganization(organizationID, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, accounts)
}

// AddAccountHandler joins a new or pre-existing account to the provided organization.
func (c *Config) AddAccountHandler(r *http.Request) *service.Response {
	account := models.Account{}
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	organizationID, err := strconv.Atoi(mux.Vars(r)["organizationID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errOrganizationIDRequired.Error()), http.StatusBadRequest, nil)
	}

	// If the account already exists, then simply add the membership to it.
	if existing, err := models.GetAccountByEmail(account.Email, c.DBClient); err != nil {
		return service.NewResponse(nil, http.StatusBadRequest, nil)
	} else if existing != nil {
		if err := existing.AddOrganization(organizationID, c.DBClient); err != nil {
			return service.NewResponse(err, http.StatusBadRequest, nil)
		}
		return service.NewResponse(nil, http.StatusOK, account)
	}

	account.Password = ""
	if err := account.CreateWithOrganization(organizationID, c.DBClient); err != nil {
		log.WithFields(log.Fields{
			"email": account.Email, "error": err.Error(), "organizationID": organizationID, "id": r.Header.Get("X-Request-ID"),
		}).Info(errAccountCreation.Error())
		return service.NewResponse(errAccountCreation, http.StatusBadRequest, map[string]string{"msg": errAccountCreation.Error()})
	}

	if err := account.SendAccountActivationEmail(c.AppDomain, c.TokenSecret, c.Clock, c.MailClient); err != nil {
		log.WithFields(log.Fields{
			"email": account.Email, "error": err.Error(), "id": r.Header.Get("X-Request-ID"),
		}).Info(errAccountActivationEmail.Error())
	}

	return service.NewResponse(nil, http.StatusOK, account)
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
