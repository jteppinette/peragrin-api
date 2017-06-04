package memberships

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"gitlab.com/peragrin/api/models"
	"gitlab.com/peragrin/api/service"
)

// ListAccountsHandler returns a response with all accounts that are
// have the provided membership.
func (c *Config) ListAccountsHandler(r *http.Request) *service.Response {
	membershipID, err := strconv.Atoi(mux.Vars(r)["membershipID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errMembershipIDRequired.Error()), http.StatusBadRequest, nil)
	}

	accounts, err := models.GetAccountsByMembership(membershipID, c.Client)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, accounts)
}
