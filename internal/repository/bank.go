package repository

import (
	"context"
	"database/sql"
	"ecomm/internal/helper/errorer"
	"ecomm/internal/model/entity"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type BankRepository interface {
	FindAll(ctx context.Context, userId int64) ([]entity.Bank, int, error)
	FindByID(ctx context.Context, id int64) (*entity.Bank, int, error)
	DeleteByID(ctx context.Context, id int64) (int, error)
	UpdateByID(ctx context.Context, ent entity.Bank) (int, error)
	Create(ctx context.Context, ent entity.Bank) (int, error)
}

func NewBankRepository(logger zerolog.Logger, db *sql.DB) BankRepository {
	return &BankRepositoryImpl{
		logger: logger,
		db:     db,
	}
}

type BankRepositoryImpl struct {
	logger zerolog.Logger
	db     *sql.DB
}

func (r *BankRepositoryImpl) FindAll(ctx context.Context, userId int64) ([]entity.Bank, int, error) {
	banks := []entity.Bank{}
	query := `SELECT 
				id,
				name, 
				account_name,
				account_number,
				user_id,
				created_at,
				updated_at
			FROM banks
			WHERE user_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, userId)

	if err != nil {
		return banks, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	for rows.Next() {
		bank := entity.Bank{}
		if err := rows.Scan(
			&bank.ID,
			&bank.Name,
			&bank.AccountName,
			&bank.AccountNumber,
			&bank.UserID,
			&bank.CreatedAt,
			&bank.UpdatedAt,
		); err != nil {
			return banks, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
		}

		banks = append(banks, bank)
	}

	return banks, http.StatusOK, nil
}

func (r *BankRepositoryImpl) DeleteByID(ctx context.Context, id int64) (int, error) {
	query := `DELETE FROM banks WHERE id = $1 RETURNING id`
	rId := 0
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&rId); err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}
	if rId == 0 {
		return http.StatusNotFound, errors.Wrap(errorer.ErrNotFound, "product not found")
	}

	return http.StatusOK, nil
}

func (r *BankRepositoryImpl) UpdateByID(ctx context.Context, ent entity.Bank) (int, error) {
	query := `
		UPDATE banks SET
			name=$1, 
			account_name=$2,
			account_number=$3,
			updated_at=$4
		Where id = $5
	`

	res, err := r.db.ExecContext(ctx, query,
		ent.Name,
		ent.AccountName,
		ent.AccountNumber,
		ent.UpdatedAt,
		ent.ID)

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

func (r *BankRepositoryImpl) Create(ctx context.Context, ent entity.Bank) (int, error) {
	// insert bank
	query := `
		Insert into banks
		(	
			name, 
			account_name,
			account_number,
			user_id,
			created_at,
			updated_at
		)
		Values($1, $2, $3, $4, $5, $6) 
	`

	_, err := r.db.ExecContext(ctx, query, ent.Name, ent.AccountName,
		ent.AccountNumber, ent.UserID, ent.CreatedAt, ent.UpdatedAt)

	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	return http.StatusOK, nil
}

func (r *BankRepositoryImpl) FindByID(ctx context.Context, id int64) (*entity.Bank, int, error) {
	bank := entity.Bank{}
	query := `
		SELECT 
			id, 
			name, 
			account_name,
			account_number,
			user_id,
			created_at,
			updated_at
		FROM banks
		WHERE banks.id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&bank.ID,
		&bank.Name,
		&bank.AccountName,
		&bank.AccountNumber,
		&bank.UserID,
		&bank.CreatedAt,
		&bank.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, errors.Wrap(errorer.ErrNotFound, errorer.ErrNotFound.Error())
		}
		return nil, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	return &bank, http.StatusOK, nil
}
