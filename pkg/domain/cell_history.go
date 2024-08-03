package domain

type CellHistory struct {
	BaseModel
	CellId      int64
	Branch      string
	Version     int64
	RequestType string
	LockKey     string
	Who         string
}

func (CellHistory) TableName() string {
	return "cell_history"
}
