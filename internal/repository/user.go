package repository

import (
	"context"
	"errors"
	"fmt"
	"marketing/internal/dto"
	"marketing/internal/infrastructure/database"

	"github.com/jackc/pgx/v5"
)

type UserRepo struct {
	postgres *database.Postgres
}

func NewUserRepo(
	postgres *database.Postgres,
) *UserRepo {
	return &UserRepo{
		postgres: postgres,
	}
}

//nolint:exhaustruct_v5
func (ur *UserRepo) GetUserByID(
	ctx context.Context,
	id uint64,
) (
	*dto.User,
	error,
) {
	pool := ur.postgres.Pool()

	u := &dto.User{}

	query := `SELECT id, name, email 
			FROM users 
			WHERE id = $1`

	err := pool.QueryRow(
		ctx,
		query,
		id,
	).Scan(&u.ID, &u.Name, &u.Email)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed scan user query %w", err)
	}

	return u, nil
}
