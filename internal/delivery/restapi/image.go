package restapi

import (
	httpHelper "ecomm/internal/helper/http"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r *Restapi) UploadImage(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, err)
	}

	imgUrl, code, err := r.service.UploadImage(c.Request().Context(), file)
	r.debugError(err)
	if err != nil {
		return httpHelper.ResponseJSONHTTP(c, code, "", nil, nil, err)
	}
	return c.JSON(code, map[string]interface{}{
		"imageUrl": imgUrl,
	})
}
