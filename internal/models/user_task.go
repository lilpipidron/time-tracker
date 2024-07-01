package models

import "time"

type UserTask struct {
	UserID    int        `gorm:"not null;index;column:user_id" json:"user_id"`
	TaskID    int        `gorm:"not null;index;column:task_id" json:"task_id"`
	StartTime time.Time  `gorm:"not null;column:start_time" json:"start_time"`
	EndTime   *time.Time `gorm:"column:end_time" json:"end_time"`
}
