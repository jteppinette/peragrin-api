package posts

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"gitlab.com/peragrin/api/models"
	"gitlab.com/peragrin/api/service"
)

// CreateHandler saves a new post to the database.
func (c *Config) CreateHandler(r *http.Request) *service.Response {
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
