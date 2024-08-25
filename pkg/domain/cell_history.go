package domain

type CellHistory struct {
	BaseModel
	Branch      string `json:"branch"`
	CellId      int64  `json:"cellId"`
	Version     int64  `json:"version"`
	RequestType string `json:"requestType"`
	LockKey     string `json:"lockKey"`
	Who         string `json:"who"`
}

func (CellHistory) TableName() string {
	return "cell_history"
}
