package repositories

import (
	"popwa/domain"
	"popwa/logs"
)

func (r popWaRepository) UpdateScore(user domain.User) error {
	result := r.db.Model(&user).Updates(user)
	if result.Error != nil {
		logs.Error("Failed to update user score")
		return result.Error
	}
	return nil
}
