package domain

import (
	"time"
)

// BaseModel base模型的定义
type BaseModel struct {
	Id          int64     `gorm:"primary_key"         json:"id" `
	UpdatedTime time.Time `gorm:"column:updated_at"   json:"updated_time" `
}

type CellModel struct {
	BaseModel
	branch        string `gorm:"type:varchar(50);column:branch" `
	FileId        int64  `gorm:"type:int;column:file_id" `
	LatestVersion int64  `gorm:"type:int;column:latest_version" `
	Status        string `gorm:"type:varchar(50);column:status" `
	LockKey       string `gorm:"type:varchar(50);column:lock_key" `
	Who           string `gorm:"type:varchar(50);column:who" `
}

func (CellModel) TableName() string {
	return "file_info"
}
