package restapi

import (
	"ecomm/internal/helper/common"
	httpHelper "ecomm/internal/helper/http"
	"ecomm/internal/model/request"
	"ecomm/internal/model/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (r *Restapi) CreateProduct(c echo.Context) error {
	req := request.Product{}
	if err := c.Bind(&req); err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, err)
	}

	req.UserID = c.Get(common.EncodedUserJwtCtxKey.ToString()).(*response.User).ID

	_, code, err := r.service.CreateProduct(c.Request().Context(), req)
	r.debugError(err)
	return httpHelper.ResponseJSONHTTP(c, code, "", nil, nil, err)
}

func (r *Restapi) DeleteProductByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, err)
	}

	code, err := r.service.DeleteProductByID(c.Request().Context(), int64(id))
	r.debugError(err)
	return httpHelper.ResponseJSONHTTP(c, code, "", nil, nil, err)
}

func (r *Restapi) PatchProductByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, err)
	}

	req := request.UpdateProduct{}
	if err := c.Bind(&req); err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, err)
	}
	req.ID = int64(id)

	_, code, err := r.service.UpdateProductByID(c.Request().Context(), req)
	r.debugError(err)
	return httpHelper.ResponseJSONHTTP(c, code, "", nil, nil, err)
}

func (r *Restapi) GetProducts(c echo.Context) error {
	req := request.GetProducts{}
	if err := c.Bind(&req); err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, err)
	}
	r.log.Debug().Msgf("GetProducts req: %+v", req)
	usr, ok := c.Get(common.EncodedUserJwtCtxKey.ToString()).(*response.User)
	if ok {
		req.UserID = usr.ID
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Offset <= 0 {
		req.Offset = 0
	}

	prd, meta, code, err := r.service.GetProducts(c.Request().Context(), req)

	r.debugError(err)
	return httpHelper.ResponseJSONHTTP(c, code, "",
		map[string]interface{}{
			"products": prd,
		}, meta, err)
}

func (r *Restapi) PatchProductStockByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, err)
	}

	req := request.UpdateProductStock{}
	if err := c.Bind(&req); err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, err)
	}
	req.ID = int64(id)

	code, err := r.service.UpdateProductStockByID(c.Request().Context(), req)
	r.debugError(err)
	return httpHelper.ResponseJSONHTTP(c, code, "", nil, nil, err)
}
