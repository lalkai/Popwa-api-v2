package repositories

import (
	"popwa/domain"
	"popwa/logs"

	"go.uber.org/zap"
)

func (r popWaRepository) GetAllUsers() (*[]domain.User, error) {
	var users []domain.User
	err := r.db.Find(&users).Error
	if err != nil {
		logs.Error("Error fetching users from the database:", zap.Error(err))
		return nil, err
	}

	return &users, nil

}
