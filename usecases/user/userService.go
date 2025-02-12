package user

import (
	"errors"
	"fit-byte/models"
	"fit-byte/types"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) UserService {
	return UserService{userRepository}
}

func (s *UserService) FindById(id string) (*models.User, error) {
	user, err := s.userRepository.FindById(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.NewError(http.StatusNotFound, "")
		}

		return nil, err
	}

	return user, nil
}

func (s *UserService) PartialUpdate(id string, payload types.UpdateUserPayload) (*models.User, error) {
	user, err := s.userRepository.PartialUpdate(id, payload)
	if err != nil {
		return nil, err
	}

	return user, nil
}