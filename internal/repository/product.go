package repository

import (
	"context"
	"database/sql"
	"ecomm/internal/helper/errorer"
	"ecomm/internal/model/entity"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type ProductRepository interface {
	FindAll(ctx context.Context) ([]entity.Product, int, error)
	FindByID(ctx context.Context, id int64) (*entity.Product, int, error)
	Create(ctx context.Context, entity entity.Product) (*entity.Product, int, error)
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &ProductRepositoryImpl{
		db: db,
	}
}

type ProductRepositoryImpl struct {
	db *sql.DB
}

func (r *ProductRepositoryImpl) FindAll(ctx context.Context) ([]entity.Product, int, error) {
	var products []entity.Product

	rows, err := r.db.QueryContext(ctx, "SELECT * FROM products")
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var product entity.Product
		err := rows.Scan(&product.ID, &product.Label, &product.LabelMask, &product.Description, &product.DescriptionMask, &product.Price, &product.UserID, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
		}
		products = append(products, product)
	}

	return products, http.StatusOK, nil
}

func (r *ProductRepositoryImpl) FindByID(ctx context.Context, id int64) (*entity.Product, int, error) {
	var product entity.Product

	row := r.db.QueryRowContext(ctx, "SELECT * FROM products WHERE id = ?", id)
	err := row.Scan(&product.ID, &product.Label, &product.LabelMask, &product.Description, &product.DescriptionMask, &product.Price, &product.UserID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, errors.Wrap(errorer.ErrNotFound, err.Error())
		}
		fmt.Println(err.Error())
		return nil, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	return &product, http.StatusOK, nil
}

func (r *ProductRepositoryImpl) Create(ctx context.Context, data entity.Product) (*entity.Product, int, error) {
	product := &entity.Product{
		Label:           data.Label,
		LabelMask:       data.LabelMask,
		Description:     data.Description,
		DescriptionMask: data.DescriptionMask,
		Price:           data.Price,
		UserID:          data.UserID,
		CreatedAt:       time.Now().UnixMilli(),
		UpdatedAt:       time.Now().UnixMilli(),
	}

	result, err := r.db.ExecContext(ctx, "INSERT INTO products (label, label_mask, description, description_mask, price, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		product.Label, product.LabelMask, product.Description, product.DescriptionMask, product.Price, product.UserID, product.CreatedAt, product.UpdatedAt)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	lastInsertID, _ := result.LastInsertId()
	product.ID = int64(lastInsertID)

	return product, http.StatusCreated, nil
}
