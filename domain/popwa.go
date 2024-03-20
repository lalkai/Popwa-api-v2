package domain

type User struct {
	UserName string `gorm:"primaryKey"`
	Score    int
}

type AddUserBody struct {
	UserName string
}

type UpdateScore struct {
	UserName string
	Score    int
}
