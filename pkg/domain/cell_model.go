package domain

import (
	"gorm.io/gorm"
	"time"
)

// BaseModel base模型的定义
type BaseModel struct {
	gorm.Model
	Id          int64     `gorm:"primaryKey;autoIncrement"            json:"id" `
	UpdatedTime time.Time `gorm:"column:updated_time"                 json:"updated_time" `
}

type CellStatus struct {
	BaseModel
	Branch        string
	CellId        int64
	LatestVersion int64
	Status        string
	LockKey       string
	Who           string
	LockTimeFrom  *time.Time
	LockTimeTo    *time.Time
}

func (CellStatus) TableName() string {
	return "cell_status"
}
