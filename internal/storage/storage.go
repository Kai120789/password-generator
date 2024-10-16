package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"password_generator/internal/dto"
	"password_generator/internal/models"
	"time"

	"go.uber.org/zap"
)

type Storage struct {
	logger   *zap.Logger
	filePath string
}

func New(l *zap.Logger, f string) *Storage {
	return &Storage{
		logger:   l,
		filePath: f,
	}
}

func (s *Storage) readUsers() ([]models.User, error) {
	file, err := os.ReadFile(s.filePath)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	// Проверка на пустой файл
	if len(file) == 0 {
		return []models.User{}, nil
	}

	var users []models.User
	err = json.Unmarshal(file, &users)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return users, nil
}

func (s *Storage) RegisterNewUser(body dto.User) (*models.User, error) {
	users, err := s.readUsers()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	id := uint(len(users) + 1)

	today := time.Now()

	user := models.User{
		ID:        id,
		Username:  body.Username,
		Password:  body.Password,
		CreatedAt: today,
	}

	users = append(users, user)

	data, err := json.Marshal(users)
	if err != nil {
		zap.S().Error("error marshalling DTO", zap.Error(err))
		return nil, err
	}

	err = os.WriteFile(s.filePath, data, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &user, nil
}

func (s *Storage) GenNewPAssword() {

}

func (s *Storage) GetAllPasswords() {

}

func (s *Storage) DeleteUserPassword() {

}
