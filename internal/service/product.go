package service

import (
	"context"
	"ecomm/internal/helper/common"
	"ecomm/internal/helper/errorer"
	"ecomm/internal/helper/validator"
	"ecomm/internal/model/entity"
	"ecomm/internal/model/request"
	"ecomm/internal/model/response"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func (s *service) GetProducts(ctx context.Context, req request.GetProducts) ([]response.Product, *common.Meta, int, error) {
	ent, meta, code, err := s.productRepo.FindAll(ctx, entity.GetAllProductFilter{
		UserOnly:       req.UserOnly,
		UserID:         req.UserID,
		Limit:          req.Limit,
		Offset:         req.Offset,
		Tags:           req.Tags,
		Condition:      req.Condition,
		ShowEmptyStock: req.ShowEmptyStock,
		MaxPrice:       req.MaxPrice,
		MinPrice:       req.MinPrice,
		SortBy:         req.SortBy,
		OrderBy:        req.OrderBy,
		Search:         req.Search,
	})

	if err != nil {
		return nil, nil, code, err
	}

	list := make([]response.Product, len(ent))
	for i, v := range ent {
		list[i] = response.Product{
			ID:            strconv.Itoa(int(v.ID)),
			Name:          v.Name,
			Price:         v.Price,
			ImageURL:      v.ImageURL,
			Stock:         v.Stock,
			UserID:        v.UserID,
			IsPurchasable: v.IsPurchasable,
			Condition:     v.Condition,
			Tags:          strings.Split(v.Tags, ","),
			PurchaseCount: v.PurchaseCount,
			CreatedAt:     v.CreatedAt,
			UpdatedAt:     v.UpdatedAt,
		}
	}

	return list, meta, http.StatusOK, nil
}

func (s *service) GetProductByID(ctx context.Context, id int64) (*response.Product, int, error) {
	ent, code, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, code, err
	}

	return &response.Product{
		ID:            strconv.Itoa(int(ent.ID)),
		Name:          ent.Name,
		Price:         ent.Price,
		ImageURL:      ent.ImageURL,
		Stock:         ent.Stock,
		UserID:        ent.UserID,
		IsPurchasable: ent.IsPurchasable,
		Condition:     ent.Condition,
		Tags:          strings.Split(ent.Tags, ","),
		PurchaseCount: 0,
		CreatedAt:     ent.CreatedAt,
		UpdatedAt:     ent.UpdatedAt,
	}, code, nil
}

func (s *service) CreateProduct(ctx context.Context, req request.Product) (*response.Product, int, error) {
	if err := validator.ValidateStruct(&req); err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}

	if req.Condition != "new" && req.Condition != "second" {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(errors.New("invalid condition")), "invalid request condition")
	}

	_, code, err := s.productRepo.Create(ctx, entity.Product{
		Name:          req.Name,
		Price:         req.Price,
		ImageURL:      req.ImageURL,
		Stock:         req.Stock,
		UserID:        req.UserID,
		IsPurchasable: req.IsPurchasable,
		Condition:     req.Condition,
		Tags:          strings.Join(req.Tags, ","),
		PurchaseCount: 0,
		CreatedAt:     time.Now().UnixMilli(),
		UpdatedAt:     time.Now().UnixMilli(),
	})

	if err != nil {
		return nil, code, err
	}

	return nil, code, nil
}

func (s *service) UpdateProductByID(ctx context.Context, req request.UpdateProduct) (*response.Product, int, error) {
	if err := validator.ValidateStruct(&req); err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}

	if req.Condition != "new" && req.Condition != "second" {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(errors.New("invalid condition")), "invalid request condition")
	}

	_, code, err := s.productRepo.UpdateByID(ctx, entity.Product{
		ID:            req.ID,
		Name:          req.Name,
		Price:         req.Price,
		ImageURL:      req.ImageURL,
		IsPurchasable: req.IsPurchasable,
		Condition:     req.Condition,
		Tags:          strings.Join(req.Tags, ","),
		UpdatedAt:     time.Now().UnixMilli(),
	})

	if err != nil {
		return nil, code, err
	}

	return nil, code, nil
}

func (s *service) UpdateProductStockByID(ctx context.Context, req request.UpdateProductStock) (int, error) {
	if err := validator.ValidateStruct(&req); err != nil {
		return http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}

	code, err := s.productRepo.UpdateStockByID(ctx, req.ID, req.Stock)

	if err != nil {
		return code, err
	}

	return code, nil
}

func (s *service) DeleteProductByID(ctx context.Context, id int64) (int, error) {
	if id == 0 {
		return http.StatusBadRequest, errors.Wrap(errors.New("invalid product id"), "invalid product id")
	}

	code, err := s.productRepo.DeleteByID(ctx, id)
	if err != nil {
		return code, err
	}
	return http.StatusOK, nil
}

func (s *service) GetProductWithSellerByID(ctx context.Context, id int64) (*response.Product, *response.SellerDetail, int, error) {
	prd, code, err := s.GetProductByID(ctx, id)
	if err != nil {
		return nil, nil, code, err
	}

	// Concurrently fetch user and total sold
	userCh := make(chan *entity.User)
	totalSoldCh := make(chan int)
	go func() {
		user, _, _ := s.userRepo.FindByID(ctx, prd.UserID)
		userCh <- user
	}()
	go func() {
		totalSold, _, _ := s.productRepo.GetTotalSoldByUserId(ctx, prd.UserID)
		totalSoldCh <- totalSold
	}()

	usr := <-userCh
	totalSold := <-totalSoldCh

	seller := response.SellerDetail{}

	if usr != nil {
		seller.ID = strconv.Itoa(int(usr.ID))
		seller.Name = usr.Name
		seller.Username = usr.Username
		seller.ProductSoldTotal = totalSold
		seller.Banks = make([]response.Bank, len(usr.Banks))

		for i, v := range usr.Banks {
			seller.Banks[i] = response.Bank{
				ID:            strconv.Itoa(int(v.ID)),
				Name:          v.Name,
				AccountName:   v.AccountName,
				AccountNumber: v.AccountNumber,
				UserID:        v.UserID,
				CreatedAt:     v.CreatedAt,
				UpdatedAt:     v.UpdatedAt,
			}
		}
	}

	return prd, &seller, code, nil
}

func (s *service) PurchaseProduct(ctx context.Context, req request.PurchaseProduct) (int, error) {
	if err := validator.ValidateStruct(&req); err != nil {
		return http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}
	bankId, _ := strconv.Atoi(req.BankAccountId)
	bank, _, err := s.bankRepo.FindByID(ctx, int64(bankId))
	if err != nil {
		return http.StatusBadRequest, err
	}
	prd, code, err := s.productRepo.FindByID(ctx, req.ProductId)

	if err != nil {
		return code, err
	}

	if prd.UserID != bank.UserID {
		return http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(errors.New("bank account not owned by seller")), "bank account not owned by seller")
	}

	code, err = s.productRepo.Purchase(ctx, req.ProductId, req.Quantity)

	return code, err
}
