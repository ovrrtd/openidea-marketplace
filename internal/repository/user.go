package repository

import (
	"context"
	"database/sql"
	"ecomm/internal/helper/errorer"
	"ecomm/internal/model/entity"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type UserRepository interface {
	Register(ctx context.Context, user entity.User) (*entity.User, int, error)
	FindByUsername(ctx context.Context, email string) (*entity.User, int, error)
	FindByID(ctx context.Context, id int64) (*entity.User, int, error)
}

func NewUserRepository(logger zerolog.Logger, db *sql.DB) UserRepository {
	return &UserRepositoryImpl{
		logger: logger,
		db:     db,
	}
}

type UserRepositoryImpl struct {
	logger zerolog.Logger
	db     *sql.DB
}

func (r *UserRepositoryImpl) Register(ctx context.Context, user entity.User) (*entity.User, int, error) {
	newUser := &entity.User{
		Name:      user.Name,
		Username:  user.Username,
		Password:  user.Password,
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}
	err := r.db.QueryRowContext(ctx, "INSERT INTO users (username, password, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		newUser.Username, newUser.Password, newUser.Name, newUser.CreatedAt, newUser.UpdatedAt).Scan(&newUser.ID)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	return newUser, http.StatusCreated, nil
}

func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (*entity.User, int, error) {
	var user entity.User

	row := r.db.QueryRowContext(ctx, "SELECT id, username, password, name, created_at, updated_at FROM users WHERE username = $1", username)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, errors.Wrap(errorer.ErrNotFound, errorer.ErrNotFound.Error())
		}
		return nil, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}
	return &user, http.StatusOK, nil
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, id int64) (*entity.User, int, error) {
	var user entity.User

	row := r.db.QueryRowContext(ctx, "SELECT id, username, password, name, created_at, updated_at FROM users WHERE id = $1", id)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, errors.Wrap(errorer.ErrNotFound, errorer.ErrNotFound.Error())
		}
		return nil, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	return &user, http.StatusOK, nil
}
