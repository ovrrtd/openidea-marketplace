package restapi

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "Ecommerce",
		Help:    "Histogram of ecomm server request duration.",
		Buckets: prometheus.LinearBuckets(1, 1, 10), // Adjust bucket sizes as needed
	}, []string{"path", "method", "status"})
)

func (r *Restapi) MakeRoute(e *echo.Echo) {
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// user
	NewRoute(e, http.MethodPost, "/v1/user/register", r.Register)
	NewRoute(e, "POST", "/v1/user/register", r.Register)
	NewRoute(e, http.MethodPost, "/v1/user/login", r.Login)
	// image
	NewRoute(e, http.MethodPost, "/v1/image", r.UploadImage, r.middleware.Authentication(true))

	// product
	NewRoute(e, http.MethodPost, "/v1/product", r.CreateProduct, r.middleware.Authentication(true))
	NewRoute(e, http.MethodGet, "/v1/product", r.GetProducts, r.middleware.Authentication(false))
	NewRoute(e, http.MethodGet, "/v1/product/:id", r.GetProductByID)
	NewRoute(e, http.MethodDelete, "/v1/product/:id", r.DeleteProductByID, r.middleware.Authentication(true), r.middleware.IsProductOwner)
	NewRoute(e, http.MethodPatch, "/v1/product/:id", r.PatchProductByID, r.middleware.Authentication(true), r.middleware.IsProductOwner)
	NewRoute(e, http.MethodPatch, "/v1/product/:id/stock", r.PatchProductStockByID, r.middleware.Authentication(true), r.middleware.IsProductOwner)
	NewRoute(e, http.MethodPost, "/v1/product/:id/buy", r.PurchaseProduct, r.middleware.Authentication(true))
	// bank
	NewRoute(e, http.MethodPost, "/v1/bank/account", r.CreateBank, r.middleware.Authentication(true))
	NewRoute(e, http.MethodGet, "/v1/bank/account", r.GetBanks, r.middleware.Authentication(false))
	NewRoute(e, http.MethodDelete, "/v1/bank/account/:id", r.DeleteBankByID, r.middleware.Authentication(true), r.middleware.IsBankOwner)
	NewRoute(e, http.MethodPatch, "/v1/bank/account/:id", r.PatchBankByID, r.middleware.Authentication(true), r.middleware.IsBankOwner)
}

func NewRoute(app *echo.Echo, method string, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) {
	app.Add(method, path, wrapHandlerWithMetrics(path, method, handler), middleware...)
}

func wrapHandlerWithMetrics(path string, method string, handler echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		startTime := time.Now()

		// Execute the actual handler and catch any errors
		err := handler(c)

		// Regardless of whether an error occurred, record the metrics
		duration := time.Since(startTime).Seconds()

		requestHistogram.WithLabelValues(path, method, strconv.Itoa(c.Response().Status)).Observe(duration)
		return err
	}
}
