package domain

import (
	"time"
)

type CellHistory struct {
	BaseModel
	CellId       int64
	Branch       string
	Version      int64
	RequestType  string
	LockKey      string
	Who          string
	LockTimeFrom *time.Time
	LockTimeTo   *time.Time
}

func (CellHistory) TableName() string {
	return "cell_history"
}
