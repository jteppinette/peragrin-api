package posts

import (
	"encoding/json"
	"net/http"

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

	// TODO: pull organization id from url, because an account might be an operator of many organizations
	// account, ok := context.GetOk(r, "account")
	// if !ok {
	// 	return service.NewResponse(errAuthenticationRequired, http.StatusUnauthorized, nil)
	// }
	// form.OrganizationID = account.(models.Account).OrganizationID
	if err := form.Save(c.Client); err != nil {
		return service.NewResponse(errors.Wrap(err, errCreatePost.Error()), http.StatusInternalServerError, nil)
	}

	return service.NewResponse(nil, http.StatusCreated, form)
}
