package repositories

import (
	"popwa/domain"
	"popwa/logs"

	"go.uber.org/zap"
)

func (r popWaRepository) GetUser(Username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("user_name = ?", Username).First(&user).Error
	if err != nil {
		logs.Error("Error fetching user from the database", zap.Error(err))
		return nil, err
	}

	return &user, nil

}
