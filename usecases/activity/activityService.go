package activity

import (
	"fit-byte/models"
	"fit-byte/utils"
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