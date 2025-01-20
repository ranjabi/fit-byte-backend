package activity

import (
	"fit-byte/constants"
	"fit-byte/models"
	"fit-byte/types"
	"fit-byte/utils"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

type ActivityService struct {
	activityRepository ActivityRepository
}

func NewActivityService(activityRepository ActivityRepository) ActivityService {
	return ActivityService{activityRepository}
}

func (s *ActivityService) CreateActivity(activity models.Activity) (*models.Activity, error) {
	newActivity, err := s.activityRepository.Save(activity)
	if err != nil {
		return nil, err
	}

	newActivity.CaloriesBurned = utils.CalculateCaloriesBurned(newActivity.ActivityType, newActivity.DurationInMinutes)

	return newActivity, nil
}

func (s *ActivityService) UpdateActivity(id string, payload types.UpdateActivityPayload) (*models.Activity, error) {
	activity, err := s.activityRepository.Update(id, payload)
	if err != nil {
		return nil, err
	}

	return activity, nil
}

func (s *ActivityService) DeleteActivity(id string) (error) {
	err := s.activityRepository.Delete(id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == constants.INVALID_INPUT_SYNTAX_TYPE_ERROR_CODE {
			return models.NewError(http.StatusNotFound, "")
		}
		
		return err
	}

	return nil
}