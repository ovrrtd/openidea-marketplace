package restapi

import "github.com/labstack/echo/v4"

func (r *Restapi) MakeRoute(e *echo.Echo) {

	v1 := e.Group("/v1")
	v1.POST("/user/register", r.Register)
	v1.POST("/user/login", r.Login)
	v1.POST("/image", r.UploadImage, r.middleware.Authentication)
}
