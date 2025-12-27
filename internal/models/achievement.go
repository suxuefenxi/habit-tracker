package models

type Achievement struct {
	ID             uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Code           string `gorm:"column:code;type:varchar(64);not null;uniqueIndex" json:"code"`
	Name           string `gorm:"column:name;type:varchar(128);not null" json:"name"`
	Description    string `gorm:"column:description;type:text" json:"description"`
	ConditionType  string `gorm:"column:condition_type;type:varchar(32);not null;index" json:"condition_type"`
	ConditionValue int    `gorm:"column:condition_value;not null" json:"condition_value"`
}

func (Achievement) TableName() string { return "achievements" }
