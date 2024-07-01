package models

type User struct {
	ID             uint   `gorm:"primaryKey;auto_increment;column:id" json:"id"`
	Name           string `gorm:"not null;column:name" json:"name"`
	Surname        string `gorm:"not null;column:surname" json:"surname"`
	Patronymic     string `gorm:"not null;column:patronymic" json:"patronymic"`
	Address        string `gorm:"not null;column:address" json:"address"`
	PassportNumber string `gorm:"not null;column:passport_number" json:"passport_number"`
}
