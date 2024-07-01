package models

import "time"

type UserTask struct {
	UserID    int        `gorm:"not null;index" json:"user_id"`
	TaskID    int        `gorm:"not null;index" json:"task_id"`
	StartTime time.Time  `gorm:"not null" json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}
