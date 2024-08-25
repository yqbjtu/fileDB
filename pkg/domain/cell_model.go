package domain

import (
	"gorm.io/gorm"
	"time"
)

// BaseModel base模型的定义

type BaseModel struct {
	gorm.Model
	Id int64 `gorm:"primaryKey;autoIncrement"            json:"id" `
}

type CellStatus struct {
	BaseModel
	Branch        string     `json:"branch"`
	CellId        int64      `json:"cellId"`
	LatestVersion int64      `json:"latestVersion"`
	Status        string     `json:"status"`
	LockKey       string     `json:"lockKey"`
	Who           string     `json:"who"`
	LockTimeFrom  *time.Time `json:"lockTimeFrom"`
	LockTimeTo    *time.Time `json:"lockTimeTo"`

	CreatedAt time.Time      `gorm:"<-:create"`      // 原样保留 gorm 设置，但在 JSON 中忽略
	UpdatedAt time.Time      `gorm:"<-:update"`      // 原样保留 gorm 设置，但在 JSON 中忽略
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 忽略该字段
}

func (CellStatus) TableName() string {
	return "cell_status"
}
