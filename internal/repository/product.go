package repository

import (
	"context"
	"database/sql"
	"ecomm/internal/helper/common"
	"ecomm/internal/helper/errorer"
	"ecomm/internal/model/entity"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type ProductRepository interface {
	FindAll(ctx context.Context, filter entity.GetAllProductFilter) ([]entity.Product, *common.Meta, int, error)
	FindByID(ctx context.Context, id int64) (*entity.Product, int, error)
	DeleteByID(ctx context.Context, id int64) (int, error)
	UpdateByID(ctx context.Context, entity entity.Product) (*entity.Product, int, error)
	UpdateStockByID(ctx context.Context, id int64, stock int) (int, error)
	Create(ctx context.Context, entity entity.Product) (*entity.Product, int, error)
	GetTotalSoldByUserId(ctx context.Context, userId int64) (int, int, error)
	Purchase(ctx context.Context, id int64, amount int) (int, error)
}

func NewProductRepository(logger zerolog.Logger, db *sql.DB) ProductRepository {
	return &ProductRepositoryImpl{
		logger: logger,
		db:     db,
	}
}

type ProductRepositoryImpl struct {
	logger zerolog.Logger
	db     *sql.DB
}

func (r *ProductRepositoryImpl) FindAll(ctx context.Context, filter entity.GetAllProductFilter) ([]entity.Product, *common.Meta, int, error) {
	var conditions []string
	var args []interface{}
	// Add conditions based on filter criteria
	argIndex := 1 // Start index for placeholder arguments
	if filter.UserOnly && filter.UserID != 0 {
		conditions = append(conditions, "user_id = $"+fmt.Sprint(argIndex))
		args = append(args, filter.UserID)
		argIndex++
	}

	if len(filter.Tags) > 0 {
		tagConditions := make([]string, len(filter.Tags))
		for i, tag := range filter.Tags {
			tagConditions[i] = "LOWER(tags) LIKE $" + fmt.Sprint(argIndex)
			args = append(args, "%"+strings.ToLower(tag)+"%")
			argIndex++
		}
		conditions = append(conditions, "("+strings.Join(tagConditions, " AND ")+")")
	}

	if filter.Condition != "" {
		conditions = append(conditions, "condition = $"+fmt.Sprint(argIndex))
		args = append(args, filter.Condition)
		argIndex++
	}

	if !filter.ShowEmptyStock {
		conditions = append(conditions, "stock > 0")
	}

	if filter.MaxPrice > 0 {
		conditions = append(conditions, "price <= $"+fmt.Sprint(argIndex))
		args = append(args, filter.MaxPrice)
		argIndex++
	}

	if filter.MinPrice > 0 {
		conditions = append(conditions, "price >= $"+fmt.Sprint(argIndex))
		args = append(args, filter.MinPrice)
		argIndex++
	}

	if filter.Search != "" {
		conditions = append(conditions, "LOWER(name) LIKE $"+fmt.Sprint(argIndex))
		args = append(args, "%"+strings.ToLower(filter.Search)+"%")
		argIndex++
	}

	// Construct the WHERE clause
	var whereClause string
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Construct the ORDER BY clause
	var orderByClause string
	if filter.SortBy != "" {
		if filter.SortBy == "date" {
			orderByClause = "ORDER BY " + "created_at"
		}

		if filter.SortBy == "price" {
			orderByClause = "ORDER BY " + "price"
		}
		orderBy := strings.ToUpper(filter.OrderBy)
		if orderBy != "" && (orderBy == "ASC" || orderBy == "DESC") {
			orderByClause += " " + orderBy
		}
	}

	// Construct the LIMIT and OFFSET clauses
	limitOffsetClause := fmt.Sprintf("LIMIT $%d ", argIndex)
	argIndex++
	limitOffsetClause += fmt.Sprintf("OFFSET $%d ", argIndex)
	argIndex++

	// Construct the final query for products
	query := `SELECT 
			id, 
			name, 
			price, 
			image_url,
			stock, 
			condition, 
			tags,
			is_purchasable,
			purchase_count, 
			user_id,
			created_at,
			updated_at
	FROM products ` + whereClause + " " + orderByClause + " " + limitOffsetClause

	// Construct the query to get total product count
	countQuery := "SELECT COUNT(*) FROM products " + whereClause

	prds := []entity.Product{}
	// Execute the main query
	argsQuery := []interface{}{}
	argsQuery = append(argsQuery, args...)
	argsQuery = append(argsQuery, filter.Limit)
	argsQuery = append(argsQuery, filter.Offset)
	rows, err := r.db.QueryContext(ctx, query, argsQuery...)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, errors.Wrap(err, "failed to execute query")
	}
	defer rows.Close()

	for rows.Next() {
		prd := entity.Product{}
		// Add more variables as needed for other columns
		err := rows.Scan(
			&prd.ID,
			&prd.Name,
			&prd.Price,
			&prd.ImageURL,
			&prd.Stock,
			&prd.Condition,
			&prd.Tags,
			&prd.IsPurchasable,
			&prd.PurchaseCount,
			&prd.UserID,
			&prd.CreatedAt,
			&prd.UpdatedAt,
		)
		if err != nil {
			return nil, nil, http.StatusInternalServerError, errors.Wrap(err, "failed to scan product")
		}
		prds = append(prds, prd)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, nil, http.StatusInternalServerError, errors.Wrap(err, "failed to iterate over rows")
	}

	// Execute the count query to get total product count
	var totalCount int
	err = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, errors.Wrap(err, "failed to get total product count")
	}

	meta := common.Meta{
		Total:  totalCount,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	}

	return prds, &meta, http.StatusOK, nil
}

func (r *ProductRepositoryImpl) FindByID(ctx context.Context, id int64) (*entity.Product, int, error) {
	prd := entity.Product{}
	usr := entity.User{}
	query := `
		SELECT 
			p.id, 
			p.name, 
			p.price, 
			p.image_url,
			p.stock, 
			p.condition, 
			p.tags,
			p.is_purchasable,
			p.purchase_count, 
			p.user_id,
			p.created_at,
			p.updated_at,
			u.name
		FROM products as p
		LEFT JOIN users as u ON p.user_id = u.id
		WHERE p.id = $1
		LIMIT 1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&prd.ID,
		&prd.Name,
		&prd.Price,
		&prd.ImageURL,
		&prd.Stock,
		&prd.Condition,
		&prd.Tags,
		&prd.IsPurchasable,
		&prd.PurchaseCount,
		&prd.UserID,
		&prd.CreatedAt,
		&prd.UpdatedAt,
		&usr.Name,
	)

	prd.User = usr

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, errors.Wrap(errorer.ErrNotFound, errorer.ErrNotFound.Error())
		}
		return nil, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	return &prd, http.StatusOK, nil
}

func (r *ProductRepositoryImpl) Create(ctx context.Context, entity entity.Product) (*entity.Product, int, error) {

	query := `
		Insert into products
		(	
			name, 
			price, 
			image_url,
			stock, 
			condition, 
			tags,
			is_purchasable,
			purchase_count, 
			user_id,
			created_at,
			updated_at
		)
		Values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
		RETURNING id;
	`

	err := r.db.QueryRowContext(ctx, query, entity.Name, entity.Price,
		entity.ImageURL, entity.Stock, entity.Condition, entity.Tags, entity.IsPurchasable,
		entity.PurchaseCount, entity.UserID, entity.CreatedAt, entity.UpdatedAt).Scan(&entity.ID)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	return &entity, http.StatusOK, nil
}

func (r *ProductRepositoryImpl) UpdateByID(ctx context.Context, entity entity.Product) (*entity.Product, int, error) {
	query := `
		UPDATE products SET
			name=$1, 
			price=$2, 
			image_url=$3,
			condition=$4, 
			tags=$5,
			is_purchasable=$6,
			updated_at=$7
		Where id = $8
	`

	res, err := r.db.ExecContext(ctx, query,
		entity.Name,
		entity.Price,
		entity.ImageURL,
		entity.Condition,
		entity.Tags,
		entity.IsPurchasable,
		entity.UpdatedAt,
		entity.ID)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}
	row, err := res.RowsAffected()

	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	if row == 0 {
		return nil, http.StatusNotFound, errors.Wrap(errorer.ErrNotFound, errorer.ErrNotFound.Error())
	}

	return &entity, http.StatusOK, nil
}

func (r *ProductRepositoryImpl) UpdateStockByID(ctx context.Context, id int64, stock int) (int, error) {
	query := `
		UPDATE products SET
			stock=$1, 
			updated_at=$2
		Where id = $3
	`

	res, err := r.db.ExecContext(ctx, query,
		stock,
		time.Now().UnixMilli(),
		id,
	)

	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}
	row, err := res.RowsAffected()

	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	if row == 0 {
		return http.StatusNotFound, errors.Wrap(errorer.ErrNotFound, errorer.ErrNotFound.Error())
	}

	return http.StatusOK, nil
}

func (r *ProductRepositoryImpl) DeleteByID(ctx context.Context, id int64) (int, error) {
	query := `DELETE FROM products WHERE id = $1 RETURNING id`
	rId := 0
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&rId); err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}
	if rId == 0 {
		return http.StatusNotFound, errors.Wrap(errorer.ErrNotFound, "product not found")
	}

	return http.StatusOK, nil
}

func (r *ProductRepositoryImpl) GetTotalSoldByUserId(ctx context.Context, userId int64) (int, int, error) {

	var total int
	query := `SELECT SUM(purchase_count) FROM products WHERE user_id = $1`
	if err := r.db.QueryRowContext(ctx, query, userId).Scan(&total); err != nil {
		return 0, 0, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	return total, http.StatusOK, nil
}

func (r *ProductRepositoryImpl) Purchase(ctx context.Context, id int64, amount int) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}
	defer tx.Rollback()

	// Acquire a row-level lock on the product row for update
	_, err = tx.ExecContext(ctx, `SELECT id FROM products WHERE id = $1 FOR UPDATE`, id)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	prd := entity.Product{}
	err = tx.QueryRowContext(ctx, `SELECT id, name, price, stock, purchase_count FROM products WHERE id = $1`, id).
		Scan(&prd.ID, &prd.Name, &prd.Price, &prd.Stock, &prd.PurchaseCount)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	if prd.Stock < amount {
		return http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(errors.New("insufficient stock")), errorer.ErrInputRequest(errors.New("insufficient stock")).Error())
	}

	prd.PurchaseCount += amount
	prd.Stock -= amount

	_, err = tx.ExecContext(ctx, `UPDATE products SET stock = $1, purchase_count = $2 WHERE id = $3`, prd.Stock, prd.PurchaseCount, prd.ID)

	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	return http.StatusOK, nil
}
