package promotions

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/jteppinette/peragrin-api/models"
	"github.com/jteppinette/peragrin-api/service"
)

// UpdateHandler updates a promotion.
func (c *Config) UpdateHandler(r *http.Request) *service.Response {
	promotion := models.Promotion{}
	if err := json.NewDecoder(r.Body).Decode(&promotion); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	id, err := strconv.Atoi(mux.Vars(r)["promotionID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errPromotionIDRequired.Error()), http.StatusBadRequest, nil)
	}
	promotion.ID = id

	if err := promotion.Save(c.Client); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, promotion)
}

// DeleteHandler deletes a promotion.
func (c *Config) DeleteHandler(r *http.Request) *service.Response {
	id, err := strconv.Atoi(mux.Vars(r)["promotionID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errPromotionIDRequired.Error()), http.StatusBadRequest, nil)
	}

	if err := models.DeletePromotion(id, c.Client); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	return service.NewResponse(nil, http.StatusOK, nil)
}

// RedeemHandler creates a account promotion relationship. This represents
// an account redeeming a promotion.
func (c *Config) RedeemHandler(r *http.Request) *service.Response {
	promotionID, err := strconv.Atoi(mux.Vars(r)["promotionID"])
	if err != nil {
		return service.NewResponse(errors.Wrap(err, errPromotionIDRequired.Error()), http.StatusBadRequest, nil)
	}

	account, ok := context.Get(r, "account").(models.Account)
	if !ok {
		return service.NewResponse(errAuthenticationRequired, http.StatusUnauthorized, nil)
	}

	redemption := &models.AccountPromotion{AccountID: account.ID, PromotionID: promotionID}

	// Does this account have the necessary membership level to redeem this promotion?
	ok, err = redemption.HasPermission(c.Client)
	if err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}
	if !ok {
		return service.NewResponse(errPromotionMembershipRequirementNotMet, http.StatusForbidden, map[string]string{"msg": errPromotionMembershipRequirementNotMet.Error()})
	}

	if err := redemption.Create(c.Client); err != nil {
		return service.NewResponse(err, http.StatusBadRequest, nil)
	}

	return service.NewResponse(nil, http.StatusOK, redemption)
}
