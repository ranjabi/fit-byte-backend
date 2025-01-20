package activity

import (
	"context"
	"fit-byte/models"
	"fit-byte/types"
	"fit-byte/utils"
	"fmt"
	"net/http"
	"strings"

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
		duration_in_minutes,
		calories_burned
	) 
	VALUES (
		@activity_type, 
		@done_at, 
		@duration_in_minutes,
		@calories_burned
	)
	RETURNING *
	`
	args := pgx.NamedArgs{
		"activity_type":       activity.ActivityType,
		"done_at":             activity.DoneAt,
		"duration_in_minutes": activity.DurationInMinutes,
		"calories_burned": utils.CalculateCaloriesBurned(activity.ActivityType, activity.DurationInMinutes),
	}

	rows, _ := r.pgConn.Query(r.ctx, query, args)
	newActivity, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Activity])
	if err != nil {
		return nil, err
	}

	return &newActivity, nil
}

func (r *ActivityRepository) GetAllActivities(offset int, limit int, activityType *string, doneAtFrom *string, doneAtTo *string, caloriesBurnedMin *string, caloriesBurnedMax *string) ([]models.Activity, error) {
	var conditions []string

	if activityType != nil {
		conditions = append(conditions, fmt.Sprintf("activity_type = '%s'", *activityType))
	}
	if doneAtFrom != nil {
		conditions = append(conditions, fmt.Sprintf("done_at >= '%s'", *doneAtFrom))
	}
	if doneAtTo != nil {
		conditions = append(conditions, fmt.Sprintf("done_at <= '%s'", *doneAtTo))
	}
	if caloriesBurnedMin != nil {
		conditions = append(conditions, fmt.Sprintf("calories_burned >= '%s'", *caloriesBurnedMin))
	}
	if caloriesBurnedMax != nil {
		conditions = append(conditions, fmt.Sprintf("calories_burned <= '%s'", *caloriesBurnedMax))
	}

	query := "SELECT * FROM activities"

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += `
	ORDER BY created_at
	LIMIT @limit
	OFFSET @offset`
	args := pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	}
	rows, err := r.pgConn.Query(r.ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("QUERY: %#v\nARGS: %#v\nROWS: %#v\n%v", query, args, rows, err.Error())
	}

	activities, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Activity])
	if err != nil {
		return nil, err
	}

	return activities, nil
}

func (r *ActivityRepository) Update(id string, payload types.UpdateActivityPayload) (*models.Activity, error) {
	if payload.DurationInMinutes != nil {
		caloriesBurned := utils.CalculateCaloriesBurned(*payload.ActivityType, *payload.DurationInMinutes)
		payload.CaloriesBurned = &caloriesBurned
	}
	query, args, err := utils.BuildPartialUpdateQuery("activities", "id", id, &payload)
	if err != nil {
		return nil, err
	}
	rows, err := r.pgConn.Query(r.ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("QUERY: %#v\nARGS: %#v\nROWS: %#v\n%v", query, args, rows, err.Error())
	}

	activity, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Activity])
	if err != nil {
		return nil, models.NewError(http.StatusNotFound, "identityId is not found")
	}

	return &activity, nil
}

func (r *ActivityRepository) Delete(id string) error {
	query := `DELETE FROM activities WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}
	commandTag, err := r.pgConn.Exec(r.ctx, query, args)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return models.NewError(http.StatusNotFound, "")
	}

	return nil
}
