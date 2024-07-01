package models

type Task struct {
	ID   int    `gorm:"primary_key;auto_increment;column:id" json:"id"`
	Name string `gorm:"not null;column:name" json:"name"`
}
