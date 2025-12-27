package models

import "time"

type Habit struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint64    `gorm:"column:user_id;not null;index" json:"user_id"`
	Name        string    `gorm:"column:name;type:varchar(128);not null" json:"name"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	TargetType  string    `gorm:"column:target_type;type:varchar(16);not null" json:"target_type"`
	TargetTimes int       `gorm:"column:target_times;not null;default:1" json:"target_times"`
	StartDate   time.Time `gorm:"column:start_date;type:date;not null" json:"start_date"`
	IsActive    bool      `gorm:"column:is_active;not null;default:true;index" json:"is_active"`
}

func (Habit) TableName() string { return "habits" }
