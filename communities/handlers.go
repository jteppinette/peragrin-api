package communities

import (
	"net/http"

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
