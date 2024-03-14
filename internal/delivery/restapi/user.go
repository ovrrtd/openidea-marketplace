package restapi

import (
	httpHelper "ecomm/internal/helper/http"
	"ecomm/internal/model/request"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r *Restapi) Register(c echo.Context) error {
	var request request.Register
	err := c.Bind(&request)
	if err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, nil)
	}

	ret, code, err := r.service.Register(c.Request().Context(), request)
	r.debugError(err)
	return httpHelper.ResponseJSONHTTP(c, code, "User registered successfully", ret, nil, err)
}

func (r *Restapi) Login(c echo.Context) error {
	var request request.Login
	err := c.Bind(&request)
	if err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, nil)
	}
	ret, code, err := r.service.Login(c.Request().Context(), request)
	r.debugError(err)
	return httpHelper.ResponseJSONHTTP(c, code, "User logged successfully", ret, nil, err)
}
