package restapi

import (
	"ecomm/internal/helper/common"
	httpHelper "ecomm/internal/helper/http"
	"ecomm/internal/model/request"
	"ecomm/internal/model/response"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (r *Restapi) CreateBank(c echo.Context) error {
	req := request.CreateBank{}
	if err := c.Bind(&req); err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, err)
	}

	req.UserID = c.Get(common.EncodedUserJwtCtxKey.ToString()).(*response.User).ID

	code, err := r.service.CreateBank(c.Request().Context(), req)
	r.debugError(err)
	return httpHelper.ResponseJSONHTTP(c, code, "", nil, nil, err)
}

func (r *Restapi) DeleteBankByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, err)
	}

	code, err := r.service.DeleteBankByID(c.Request().Context(), int64(id))
	r.debugError(err)
	return httpHelper.ResponseJSONHTTP(c, code, "", nil, nil, err)
}

func (r *Restapi) PatchBankByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, err)
	}

	req := request.UpdateBank{}
	if err := c.Bind(&req); err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, err)
	}
	req.ID = int64(id)

	code, err := r.service.UpdateBankByID(c.Request().Context(), req)
	r.debugError(err)
	return httpHelper.ResponseJSONHTTP(c, code, "", nil, nil, err)
}

func (r *Restapi) GetBanks(c echo.Context) error {
	usr, ok := c.Get(common.EncodedUserJwtCtxKey.ToString()).(*response.User)
	if !ok {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, errors.New("invalid user"))
	}
	banks, code, err := r.service.GetBanks(c.Request().Context(), usr.ID)

	r.debugError(err)
	return httpHelper.ResponseJSONHTTP(c, code, "",
		banks, nil, err)
}
