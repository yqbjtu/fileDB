package domain

import "github.com/golang/protobuf/ptypes/duration"

type CellBase struct {
	CellId    string `json:"cellId"`
	Version   int64  `json:"version"`
	Namespace string `json:"namespace"`
}

type AddVersionReq struct {
	CellBase
	LockKey string `json:"LockKey"`
	Comment string `json:"comment"`
}

type LockReq struct {
	CellId       string            `json:"cellId"`
	Namespace    string            `json:"namespace"`
	LockKey      string            `json:"LockKey"`
	lockDuration duration.Duration `json:"lockDuration"`
}

/*
示例而已，因此字段只有几个
*/
type User struct {
	UserId int64
	AddVersionReq
	//etc ..
}

type CellStatus struct {
	CellId        string
	LatestVersion int64
	Namespace     string
	user          string
	addTime       string
	IsLocked      string
	LockKey       string
	LockTimeFrom  string
	LockTimeTo    string
}
