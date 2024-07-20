package domain

import "github.com/golang/protobuf/ptypes/duration"

type CellBase struct {
	CellId  string `json:"cellId"`
	Version int64  `json:"version"`
	Branch  string `json:"Branch"`
}

type AddVersionReq struct {
	CellBase
	LockKey string `json:"LockKey"`
	Comment string `json:"comment"`
}

type LockReq struct {
	CellId       string            `json:"cellId"`
	Branch       string            `json:"Branch"`
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
