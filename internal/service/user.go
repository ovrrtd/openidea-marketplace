package service

import (
	"context"
	"ecomm/internal/helper/common"
	"ecomm/internal/helper/errorer"
	"ecomm/internal/helper/jwt"
	"ecomm/internal/helper/validator"
	"ecomm/internal/model/entity"
	"ecomm/internal/model/request"
	"ecomm/internal/model/response"
	"net/http"
	"time"

	jwtV5 "github.com/golang-jwt/jwt/v5"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Register to register a new user by email and password
func (s *service) Register(ctx context.Context, payload request.Register) (*response.Login, int, error) {
	err := validator.ValidateStruct(&payload)

	if err != nil {
		return nil, http.StatusBadRequest, errorer.ErrInputRequest(err)
	}

	exist, code, err := s.userRepo.FindByUsername(ctx, payload.Username)

	if err != nil && code != http.StatusNotFound {
		return nil, code, err
	}
	if exist != nil {
		return nil, http.StatusConflict, errors.Wrap(errorer.ErrEmailExist, errorer.ErrEmailExist.Error())
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, err.Error())
	}

	user, code, err := s.userRepo.Register(ctx, entity.User{
		Username: payload.Username,
		Name:     payload.Name,
		Password: string(hashedPassword),
	})

	if err != nil {
		return nil, code, err
	}

	// TODO: generate access token
	userClaims := common.UserClaims{
		Id: user.ID,
		RegisteredClaims: jwtV5.RegisteredClaims{
			IssuedAt:  jwtV5.NewNumericDate(time.Now()),
			ExpiresAt: jwtV5.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		},
	}
	tokenString, err := jwt.GenerateJwt(userClaims)

	if err != nil {
		return nil, code, errors.Wrap(err, err.Error())
	}

	return &response.Login{
		Name:        user.Name,
		Username:    user.Username,
		AccessToken: tokenString,
	}, code, nil
}

func (s *service) Login(ctx context.Context, payload request.Login) (*response.Login, int, error) {
	err := validator.ValidateStruct(&payload)

	if err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}

	user, code, err := s.userRepo.FindByUsername(ctx, payload.Username)
	if err != nil {
		return nil, code, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(err, err.Error())
	}

	// TODO: generate access token
	userClaims := common.UserClaims{
		Id: user.ID,
		RegisteredClaims: jwtV5.RegisteredClaims{
			IssuedAt:  jwtV5.NewNumericDate(time.Now()),
			ExpiresAt: jwtV5.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		},
	}
	tokenString, err := jwt.GenerateJwt(userClaims)

	if err != nil {
		return nil, code, errors.Wrap(err, err.Error())
	}

	return &response.Login{
		Name:        user.Name,
		Username:    user.Username,
		AccessToken: tokenString,
	}, code, nil
}

func (s *service) GetUserByID(ctx context.Context, id int64) (*response.User, int, error) {
	user, code, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, code, err
	}
	return &response.User{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, code, nil
}
