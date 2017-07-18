package accounts

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"gitlab.com/peragrin/api/models"
	"gitlab.com/peragrin/api/service"
)

// ForgotPasswordHandler generates a token that can be used to reset the password
// of the account with the provided email address.
func (c *Config) ForgotPasswordHandler(r *http.Request) *service.Response {
	id, err := strconv.Atoi(mux.Vars(r)["accountID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errAccountIDRequired.Error()), http.StatusBadRequest, nil)
	}

	account, err := models.GetAccountByID(id, c.DBClient)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	if account == nil {
		return service.NewResponse(errAccountNotFound, http.StatusNotFound, nil)
	}

	if err := account.SendResetPasswordEmail(c.AppDomain, c.TokenSecret, c.Clock, c.MailClient); err != nil {
		return service.NewResponse(err, http.StatusInternalServerError, nil)
	}
	return service.NewResponse(nil, http.StatusOK, nil)
}
