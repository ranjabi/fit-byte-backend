package auth

import (
	"fit-byte/constants"
	"fit-byte/models"
	"fit-byte/usecases/user"
	"fit-byte/utils"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

type AuthService struct {
	userRepository user.UserRepository
}

func NewAuthService(userRepository user.UserRepository) AuthService {
	return AuthService{userRepository}
}

func (s *AuthService) CreateUser(user models.User) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	newUser, err := s.userRepository.Save(user)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == constants.UNIQUE_VIOLATION_ERROR_CODE {
			return nil, models.NewError(http.StatusConflict, "Email is already taken")
		}
		return nil, err
	}

	token, err := utils.CreateClaims(newUser)
	if err != nil {
		return nil, err
	}
	newUser.Token = token

	return newUser, nil
}