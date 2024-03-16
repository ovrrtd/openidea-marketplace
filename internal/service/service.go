package service

import (
	"context"
	"ecomm/internal/helper/common"
	"ecomm/internal/model/request"
	"ecomm/internal/model/response"
	"ecomm/internal/repository"
	"mime/multipart"

	"github.com/rs/zerolog"
)

type Service interface {
	// Product
	GetProducts(ctx context.Context, req request.GetProducts) ([]response.Product, *common.Meta, int, error)
	GetProductByID(ctx context.Context, id int64) (*response.Product, int, error)
	GetProductWithSellerByID(ctx context.Context, id int64) (*response.Product, *response.SellerDetail, int, error)
	DeleteProductByID(ctx context.Context, id int64) (int, error)
	CreateProduct(ctx context.Context, req request.Product) (*response.Product, int, error)
	UpdateProductByID(ctx context.Context, req request.UpdateProduct) (*response.Product, int, error)
	UpdateProductStockByID(ctx context.Context, req request.UpdateProductStock) (int, error)
	PurchaseProduct(ctx context.Context, req request.PurchaseProduct) (int, error)
	// User
	Register(ctx context.Context, payload request.Register) (*response.Login, int, error)
	Login(ctx context.Context, payload request.Login) (*response.Login, int, error)
	GetUserByID(ctx context.Context, id int64) (*response.User, int, error)
	// s3
	UploadImage(ctx context.Context, file *multipart.FileHeader) (string, int, error)
	// bank
	GetBanks(ctx context.Context, userId int64) ([]response.Bank, int, error)
	GetBankByID(ctx context.Context, id int64) (*response.Bank, int, error)
	DeleteBankByID(ctx context.Context, id int64) (int, error)
	UpdateBankByID(ctx context.Context, ent request.UpdateBank) (int, error)
	CreateBank(ctx context.Context, ent request.CreateBank) (int, error)
}

type Config struct {
	Salt      int
	JwtSecret string
}

type service struct {
	cfg         Config
	log         zerolog.Logger
	productRepo repository.ProductRepository
	userRepo    repository.UserRepository
	s3Repo      repository.S3Repository
	bankRepo    repository.BankRepository
}

func New(cfg Config, logger zerolog.Logger, productRepo repository.ProductRepository, userRepo repository.UserRepository, s3Repo repository.S3Repository, bankRepo repository.BankRepository) Service {
	return &service{
		cfg:         cfg,
		log:         logger,
		productRepo: productRepo,
		userRepo:    userRepo,
		s3Repo:      s3Repo,
		bankRepo:    bankRepo,
	}
}
