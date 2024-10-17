package service

import (
	"password_generator/internal/dto"
	"password_generator/internal/models"
)

type Service struct {
	storage Storager
}

type Storager interface {
	RegisterNewUser(body dto.User) (*models.User, error)
	GenNewPassword()
	GetAllPasswords(username string) (*[]models.User, error)
	DeleteUserPassword(username string, password string) error
}

func New(s Storager) *Service {
	return &Service{
		storage: s,
	}
}

func (s *Service) RegisterNewUser(body dto.User) (*models.User, error) {
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
	passwords, err := s.storage.GetAllPasswords(username)
	if err != nil {
		return nil, err
	}

	return passwords, nil
}

func (s *Service) DeleteUserPassword(username string, password string) error {
	err := s.storage.DeleteUserPassword(username, password)
	if err != nil {
		return err
	}

	return nil
}
