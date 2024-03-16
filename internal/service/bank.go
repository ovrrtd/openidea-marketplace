package service

import (
	"context"
	"ecomm/internal/helper/errorer"
	"ecomm/internal/helper/validator"
	"ecomm/internal/model/entity"
	"ecomm/internal/model/request"
	"ecomm/internal/model/response"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

func (s *service) GetBanks(ctx context.Context, userId int64) ([]response.Bank, int, error) {
	ent, code, err := s.bankRepo.FindAll(ctx, userId)

	if err != nil {
		return nil, code, err
	}

	list := make([]response.Bank, len(ent))
	for i, v := range ent {
		list[i] = response.Bank{
			ID:            strconv.Itoa(int(v.ID)),
			Name:          v.Name,
			AccountName:   v.AccountName,
			AccountNumber: v.AccountNumber,
			UserID:        v.UserID,
			CreatedAt:     v.CreatedAt,
			UpdatedAt:     v.UpdatedAt,
		}
	}

	return list, http.StatusOK, nil
}
func (s *service) GetBankByID(ctx context.Context, id int64) (*response.Bank, int, error) {
	ent, code, err := s.bankRepo.FindByID(ctx, id)
	if err != nil {
		return nil, code, err
	}

	return &response.Bank{
		ID:            strconv.Itoa(int(ent.ID)),
		Name:          ent.Name,
		AccountName:   ent.AccountName,
		AccountNumber: ent.AccountNumber,
		UserID:        ent.UserID,
		CreatedAt:     ent.CreatedAt,
		UpdatedAt:     ent.UpdatedAt,
	}, code, nil
}
func (s *service) DeleteBankByID(ctx context.Context, id int64) (int, error) {
	if id == 0 {
		return http.StatusBadRequest, errors.Wrap(errors.New("invalid bank id"), "invalid bank id")
	}

	code, err := s.bankRepo.DeleteByID(ctx, id)
	if err != nil {
		return code, err
	}
	return http.StatusOK, nil
}
func (s *service) UpdateBankByID(ctx context.Context, req request.UpdateBank) (int, error) {
	if err := validator.ValidateStruct(&req); err != nil {
		return http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}

	code, err := s.bankRepo.UpdateByID(ctx, entity.Bank{
		ID:            req.ID,
		Name:          req.Name,
		AccountName:   req.AccountName,
		AccountNumber: req.AccountNumber,
		UpdatedAt:     time.Now().UnixMilli(),
	})

	if err != nil {
		return code, err
	}

	return code, nil
}
func (s *service) CreateBank(ctx context.Context, req request.CreateBank) (int, error) {
	if err := validator.ValidateStruct(&req); err != nil {
		return http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}

	code, err := s.bankRepo.Create(ctx, entity.Bank{
		Name:          req.Name,
		AccountName:   req.AccountName,
		AccountNumber: req.AccountNumber,
		UserID:        req.UserID,
		CreatedAt:     time.Now().UnixMilli(),
		UpdatedAt:     time.Now().UnixMilli(),
	})

	if err != nil {
		return code, err
	}

	return code, nil
}
