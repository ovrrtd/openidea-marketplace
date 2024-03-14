package service

import (
	"context"
	"ecomm/internal/helper/errorer"
	"ecomm/internal/helper/validator"
	"ecomm/internal/model/entity"
	"ecomm/internal/model/request"
	"ecomm/internal/model/response"
	"math/rand"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

func (s *service) GetProducts(ctx context.Context) ([]response.Product, int, error) {
	ent, code, err := s.productRepo.FindAll(ctx)

	if err != nil {
		return nil, code, err
	}

	list := make([]response.Product, len(ent))
	for i, v := range ent {
		list[i] = response.Product{
			ID:              v.ID,
			Label:           v.Label,
			LabelMask:       v.LabelMask,
			Description:     v.Description,
			DescriptionMask: v.DescriptionMask,
			Price:           v.Price,
			UserID:          v.UserID,
			CreatedAt:       v.CreatedAt,
			UpdatedAt:       v.UpdatedAt,
		}
	}

	return list, code, nil
}

func (s *service) GetProductByID(ctx context.Context, id int64) (*response.Product, int, error) {
	ent, code, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, code, err
	}

	return &response.Product{
		ID:              ent.ID,
		Label:           ent.Label,
		LabelMask:       ent.LabelMask,
		Description:     ent.Description,
		DescriptionMask: ent.DescriptionMask,
		Price:           ent.Price,
		UserID:          ent.UserID,
		CreatedAt:       ent.CreatedAt,
		UpdatedAt:       ent.UpdatedAt,
	}, code, nil
}

func (s *service) CreateProduct(ctx context.Context, req request.Product) (*response.Product, int, error) {
	if err := validator.ValidateStruct(&req); err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}

	en, code, err := s.productRepo.Create(ctx, entity.Product{
		Label:           req.Label,
		LabelMask:       s.maskText(req.Label),
		Description:     req.Description,
		DescriptionMask: s.maskText(req.Description), // todo masking the description to be used as an input in instagram post caption
		Price:           req.Price,
		UserID:          req.UserID,
	})

	if err != nil {
		return nil, code, err
	}

	return &response.Product{
		ID:              en.ID,
		Label:           en.Label,
		LabelMask:       en.LabelMask,
		Description:     en.Description,
		DescriptionMask: en.DescriptionMask,
		Price:           en.Price,
		UserID:          en.UserID,
		CreatedAt:       en.CreatedAt,
		UpdatedAt:       en.UpdatedAt,
	}, code, nil
}

func (s *service) maskText(str string) string {
	temp := ""
	output := ""

	for i, v := range str {
		if i == len(str)-1 || v == ' ' {
			output += s.maskRandomLetters(temp, len(temp)/3)
			output += string(v)
			temp = ""
			continue
		}
		temp += string(v)
	}

	return output
}

func (s *service) maskRandomLetters(word string, numToCensor int) string {
	rand.Seed(time.Now().UnixNano())

	if len(word) == 0 || numToCensor <= 0 || numToCensor > len(word) {
		return word
	}

	// Convert the word string to a slice of runes to modify individual characters
	wordSlice := []rune(word)

	// Create a list of unique random indices to censor
	indicesToCensor := s.generateRandomIndices(len(word), numToCensor)

	// Replace the characters at the selected indices with '*'
	for _, index := range indicesToCensor {
		wordSlice[index] = '*'
	}

	// Convert the modified rune slice back to a string
	censoredWord := string(wordSlice)

	return censoredWord
}

func (s *service) generateRandomIndices(maxIndex, count int) []int {
	randIndices := make([]int, count)
	for i := 0; i < count; i++ {
		randIndices[i] = rand.Intn(maxIndex)
	}
	return randIndices
}
