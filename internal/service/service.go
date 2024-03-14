package service

import (
	"context"
	"ecomm/internal/model/request"
	"ecomm/internal/model/response"
	"ecomm/internal/repository"
	"mime/multipart"

	"github.com/rs/zerolog"
)

type Service interface {
	// Product
	GetProducts(ctx context.Context) ([]response.Product, int, error)
	GetProductByID(ctx context.Context, id int64) (*response.Product, int, error)
	CreateProduct(ctx context.Context, req request.Product) (*response.Product, int, error)
	// User
	Register(ctx context.Context, payload request.Register) (*response.Login, int, error)
	Login(ctx context.Context, payload request.Login) (*response.Login, int, error)
	GetUserByID(ctx context.Context, id int64) (*response.User, int, error)
	// s3
	UploadImage(ctx context.Context, file *multipart.FileHeader) (string, int, error)
}

type service struct {
	log         zerolog.Logger
	productRepo repository.ProductRepository
	userRepo    repository.UserRepository
	s3Repo      repository.S3Repository
}

func New(logger zerolog.Logger, productRepo repository.ProductRepository, userRepo repository.UserRepository, s3Repo repository.S3Repository) Service {
	return &service{
		log:         logger,
		productRepo: productRepo,
		userRepo:    userRepo,
		s3Repo:      s3Repo,
	}
}
