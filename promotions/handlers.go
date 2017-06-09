package promotions

import (
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"gitlab.com/peragrin/api/models"
	"gitlab.com/peragrin/api/service"
)

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
