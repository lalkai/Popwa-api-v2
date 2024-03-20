package repositories

import (
	"popwa/domain"
	"popwa/logs"
)

func (r popWaRepository) AddUser(user domain.User) error {
	result := r.db.Create(&user)
	if result.Error != nil {
		logs.Error("Failed to create user")
		return result.Error
	}
	return nil
}
