package user

import (
	"context"
	"fit-byte/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	ctx    context.Context
	pgConn *pgxpool.Pool
}

func NewUserRepository(ctx context.Context, pgConn *pgxpool.Pool) UserRepository {
	return UserRepository{ctx, pgConn}
}

func (r *UserRepository) Save(user models.User) (*models.User, error) {
	query := `
	INSERT INTO users (
		email,
		password
	) 
	VALUES (
		@email,
		@password
	)
	RETURNING *
	`
	args := pgx.NamedArgs{
		"email":    user.Email,
		"password": user.Password,
	}

	rows, _ := r.pgConn.Query(r.ctx, query, args)
	newUser, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}