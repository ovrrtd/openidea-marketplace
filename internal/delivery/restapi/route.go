package restapi

import "github.com/labstack/echo/v4"

func (r *Restapi) MakeRoute(e *echo.Echo) {

	v1 := e.Group("/v1")
	// user
	v1.POST("/user/register", r.Register)
	v1.POST("/user/login", r.Login)
	// image
	v1.POST("/image", r.UploadImage, r.middleware.Authentication(true))

	// product
	v1.POST("/product", r.CreateProduct, r.middleware.Authentication(true))
	v1.GET("/product", r.GetProducts, r.middleware.Authentication(false))
	v1.GET("/product/:id", r.GetProductByID)
	v1.DELETE("/product/:id", r.DeleteProductByID, r.middleware.Authentication(true), r.middleware.IsProductOwner)
	v1.PATCH("/product/:id", r.PatchProductByID, r.middleware.Authentication(true), r.middleware.IsProductOwner)
	v1.PATCH("/product/:id/stock", r.PatchProductStockByID, r.middleware.Authentication(true), r.middleware.IsProductOwner)
	v1.POST("/product/:id/buy", r.PurchaseProduct, r.middleware.Authentication(true))
	// bank
	v1.POST("/bank/account", r.CreateBank, r.middleware.Authentication(true))
	v1.GET("/bank/account", r.GetBanks, r.middleware.Authentication(false))
	v1.DELETE("/bank/account/:id", r.DeleteBankByID, r.middleware.Authentication(true), r.middleware.IsBankOwner)
	v1.PATCH("/bank/account/:id", r.PatchBankByID, r.middleware.Authentication(true), r.middleware.IsBankOwner)
}
