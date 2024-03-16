package middleware

import (
	"ecomm/internal/helper/common"
	"ecomm/internal/helper/errorer"
	httpHelper "ecomm/internal/helper/http"
	"ecomm/internal/helper/jwt"
	"ecomm/internal/model/response"
	"ecomm/internal/service"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type middleware struct {
	logger  zerolog.Logger
	service service.Service
}

type Middleware interface {
	Authentication(isThrowError bool) func(next echo.HandlerFunc) echo.HandlerFunc
	IsProductOwner(next echo.HandlerFunc) echo.HandlerFunc
	IsBankOwner(next echo.HandlerFunc) echo.HandlerFunc
}

func New(logger zerolog.Logger, service service.Service) Middleware {
	return &middleware{
		logger:  logger,
		service: service,
	}
}

func (m *middleware) Authentication(isThrowError bool) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			m.logger.Info().Msg("Authentication")
			token := httpHelper.GetJWTFromRequest(c.Request())

			if token == "" && isThrowError {
				return httpHelper.ResponseJSONHTTP(c, http.StatusUnauthorized, "", nil, nil, errorer.ErrUnauthorized)
			}

			if token != "" {
				claims := &common.UserClaims{}
				err := jwt.VerifyJwt(token, claims, os.Getenv("JWT_SECRET"))
				if err != nil {
					return httpHelper.ResponseJSONHTTP(c, http.StatusForbidden, "", nil, nil, errorer.ErrForbidden)
				}

				usr, code, err := m.service.GetUserByID(c.Request().Context(), claims.Id)
				if err != nil {
					return httpHelper.ResponseJSONHTTP(c, code, "", nil, nil, err)
				}
				c.Set(common.EncodedUserJwtCtxKey.ToString(), usr)
			}

			return next(c)
		}
	}
}

func (m *middleware) IsProductOwner(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		usr := c.Get(common.EncodedUserJwtCtxKey.ToString()).(*response.User)
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, err)
		}

		prd, code, err := m.service.GetProductByID(c.Request().Context(), int64(id))
		if err != nil {
			m.logger.Debug().Stack().Err(err).Send()
			return httpHelper.ResponseJSONHTTP(c, code, "", nil, nil, err)
		}

		if prd.UserID != usr.ID {
			return httpHelper.ResponseJSONHTTP(c, http.StatusForbidden, "", nil, nil, errorer.ErrForbidden)
		}
		m.logger.Debug().Msg("IsProductOwner check done")
		return next(c)
	}
}

func (m *middleware) IsBankOwner(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		usr := c.Get(common.EncodedUserJwtCtxKey.ToString()).(*response.User)
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, err)
		}

		prd, code, err := m.service.GetBankByID(c.Request().Context(), int64(id))
		if err != nil {
			m.logger.Debug().Stack().Err(err).Send()
			return httpHelper.ResponseJSONHTTP(c, code, "", nil, nil, err)
		}

		if prd.UserID != usr.ID {
			return httpHelper.ResponseJSONHTTP(c, http.StatusForbidden, "", nil, nil, errorer.ErrForbidden)
		}
		return next(c)
	}
}
