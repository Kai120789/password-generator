package service

import (
	"passwordgenerator/internal/dto"
	"passwordgenerator/internal/models"
	"passwordgenerator/internal/utils/genpass"
)

type Service struct {
	storage Storager
}

type Storager interface {
	RegisterNewUser(body dto.User) (*models.User, error)
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

func (s *Service) GenNewPassword(username string) (*models.User, error) {
	password := genpass.GeneratePassword()

	dto := dto.User{
		Username: username,
		Password: password,
	}

	user, err := s.storage.RegisterNewUser(dto)
	if err != nil {
		return nil, err
	}

	return user, nil
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
