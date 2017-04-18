package users

import (
	"net/http"

	"github.com/pkg/errors"

	"gitlab.com/peragrin/api/models"
	"gitlab.com/peragrin/api/service"
)

// ListHandler writes all users in the database to the provided response writer.
func (c *Config) ListHandler(r *http.Request) *service.Response {
	v, err := models.ListUsers(c.Client)
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errListUsers.Error()), http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, v)
}
