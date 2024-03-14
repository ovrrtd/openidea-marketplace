package middleware

import (
	"ecomm/internal/helper/common"
	"ecomm/internal/helper/errorer"
	httpHelper "ecomm/internal/helper/http"
	"ecomm/internal/helper/jwt"
	"ecomm/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type middleware struct {
	logger  zerolog.Logger
	service service.Service
}

type Middleware interface {
	Authentication(next echo.HandlerFunc) echo.HandlerFunc
}

func New(logger zerolog.Logger, service service.Service) Middleware {
	return &middleware{
		logger:  logger,
		service: service,
	}
}

func (m *middleware) Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		m.logger.Info().Msg("Authentication")
		token := httpHelper.GetJWTFromRequest(c.Request())

		if token == "" {
			return httpHelper.ResponseJSONHTTP(c, http.StatusUnauthorized, "", nil, nil, errorer.ErrUnauthorized)
		}
		m.logger.Info().Msgf("Token: %s", token)
		// TODO: validate token
		claims := &common.UserClaims{}
		err := jwt.VerifyJwt(token, claims)
		if err != nil {
			return httpHelper.ResponseJSONHTTP(c, http.StatusForbidden, "", nil, nil, errorer.ErrForbidden)
		}
		m.logger.Info().Msgf("Claims: %v", claims)

		usr, code, err := m.service.GetUserByID(c.Request().Context(), claims.Id)
		if err != nil {
			return httpHelper.ResponseJSONHTTP(c, code, "", nil, nil, err)
		}
		m.logger.Info().Msgf("User: %v", usr)
		c.Set(common.EncodedUserJwtCtxKey.ToString(), usr)
		m.logger.Info().Msg("Authentication done")
		return next(c)
	}
}
