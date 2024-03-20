package repositories

import (
	"popwa/domain"
	"gorm.io/gorm"
)

type popWaRepository struct {
	db     *gorm.DB
}

func NewPopWaRepository(db *gorm.DB) *popWaRepository {
	db.AutoMigrate(&domain.User{})
	//db.AutoMigrate(&domain.Position{})
	return &popWaRepository{
		db:     db,
	}
}
