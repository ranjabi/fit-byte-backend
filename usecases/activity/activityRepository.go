package activity

import (
	"context"
	"fit-byte/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ActivityRepository struct {
	ctx    context.Context
	pgConn *pgxpool.Pool
}

func NewActivityRepository(ctx context.Context, pgConn *pgxpool.Pool) ActivityRepository {
	return ActivityRepository{ctx, pgConn}
}

func (r *ActivityRepository) Save(activity models.Activity) (*models.Activity, error) {
	query := `
	INSERT INTO activities (
		activity_type, 
		done_at, 
		duration_in_minutes
	) 
	VALUES (
		@activity_type, 
		@done_at, 
		@duration_in_minutes
	)
	RETURNING *
	`
	args := pgx.NamedArgs{
		"activity_type": activity.ActivityType,
		"done_at": activity.DoneAt,
		"duration_in_minutes": activity.DurationInMinutes,
	}

	rows, _ := r.pgConn.Query(r.ctx, query, args)
	newActivity, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Activity])
	if err != nil {
		return nil, err
	}

	return &newActivity, nil
}
