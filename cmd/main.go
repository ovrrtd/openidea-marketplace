package main

import (
	mw "ecomm/internal/delivery/middleware"
	"ecomm/internal/delivery/restapi"
	"ecomm/internal/repository"
	"ecomm/internal/service"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	_ "github.com/lib/pq"
)

func main() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	logger := zerolog.New(os.Stdout)

	// db, err := newMongoDB(ConfigMongoDB{Host: cfg.DB.Host})
	db, err := newDBDefaultSql()
	if err != nil {
		logger.Info().Msg(fmt.Sprintf("Postgres connection error: %s", err.Error()))
		return
	}
	logger.Info().Msg(fmt.Sprintf("Postgres connected: %s", DB_HOST))
	err = db.Ping()
	if err != nil {
		logger.Info().Msg(fmt.Sprintf("Postgres ping error: %s", err.Error()))
		return
	}
	defer db.Close()

	// repository init
	productRepo := repository.NewProductRepository(logger, db)
	userRepo := repository.NewUserRepository(logger, db)
	bankRepo := repository.NewBankRepository(logger, db)
	s3Repo := repository.NewS3Repository(logger)
	salt, err := strconv.Atoi(os.Getenv("BCRYPT_SALT"))
	if err != nil {
		salt = 8
	}
	// service registry
	service := service.New(
		service.Config{Salt: salt, JwtSecret: os.Getenv("JWT_SECRET")},
		logger, productRepo, userRepo, s3Repo, bankRepo)

	// middleware init
	md := mw.New(logger, service)

	// restapi init
	rest := restapi.New(logger, md, service)

	// echo server
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("URI", v.URI).
				Str("method", c.Request().Method).
				Int("status", v.Status).
				Msg("request")
			return nil
		},
	}))

	// add restapi route
	rest.MakeRoute(e)

	errs := make(chan error)
	go func() {
		logger.Log().Msg(fmt.Sprintf("start server on port %s", APP_PORT))
		errs <- e.Start(fmt.Sprintf(":%s", APP_PORT))
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	<-errs
}
