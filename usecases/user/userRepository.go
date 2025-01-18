package user

import (
	"context"
	"fit-byte/models"
	"fit-byte/types"
	"fit-byte/utils"
	"fmt"

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

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	query := `SELECT * FROM users WHERE LOWER(email) = LOWER(@email)`
	args := pgx.NamedArgs{
		"email": email,
	}

	rows, _ := r.pgConn.Query(r.ctx, query, args)
	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindById(id string) (*models.User, error) {
	query := `SELECT * FROM users WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}

	rows, _ := r.pgConn.Query(r.ctx, query, args)
	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) PartialUpdate(id string, payload types.UpdateUserPayload) (*models.User, error) {
	query, args, err := utils.BuildPartialUpdateQuery("users", "id", id, &payload)
	if err != nil {
		return nil, err
	}
	
	rows, err := r.pgConn.Query(r.ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("QUERY: %#v\nARGS: %#v\nROWS: %#v\n%v", query, args, rows, err.Error())
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return nil, err
	}

	return &user, nil
}