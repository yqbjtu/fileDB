package domain

import (
	"time"
)

// BaseModel base模型的定义
type BaseModel struct {
	Id          int64     `gorm:"primaryKey;autoIncrement"         json:"id" `
	UpdatedTime time.Time `gorm:"column:updated_at"                 json:"updated_time" `
}

type CellStatus struct {
	BaseModel
	Branch        string
	CellId        int64
	LatestVersion int64
	Status        string
	LockKey       string
	Who           string
	LockTimeFrom  string
	LockTimeTo    string
}

func (CellStatus) TableName() string {
	return "cell_status"
}
