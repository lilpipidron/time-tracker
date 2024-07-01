package models

type Task struct {
	ID   int    `gorm:"primary_key" json:"id"`
	Name string `gorm:"not null" json:"name"`
}
