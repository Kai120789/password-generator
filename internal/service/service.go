package service

import (
	"password_generator/internal/dto"
	"password_generator/internal/models"
	hash "password_generator/internal/utils/passhash"
)

type Service struct {
	storage Storager
}

type Storager interface {
	RegisterNewUser(body dto.User) (*models.User, error)
	GenNewPAssword()
	GetAllPasswords()
	DeleteUserPassword()
}

func New(s Storager) *Service {
	return &Service{
		storage: s,
	}
}

func (s *Service) RegisterNewUser(body dto.User) (*models.User, error) {
	passwordHash, err := hash.HashPassword(body.Password)
	if err != nil {
		return nil, err
	}

	body.Password = passwordHash

	user, err := s.storage.RegisterNewUser(body)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GenNewPassword(body dto.User) (*models.User, error) {
	return nil, nil
}

func (s *Service) GetAllPasswords(username string) (*[]models.User, error) {
	return nil, nil
}

func (s *Service) DeleteUserPassword(username string, password string) error {
	return nil
}
