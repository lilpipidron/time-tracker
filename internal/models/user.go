package models

type User struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	Name           string `gorm:"not null" json:"name"`
	Surname        string `gorm:"not null" json:"surname"`
	Patronymic     string `gorm:"not null" json:"patronymic"`
	Address        string `gorm:"not null" json:"address"`
	PassportNumber string `gorm:"not null" json:"passport_number"`
}
