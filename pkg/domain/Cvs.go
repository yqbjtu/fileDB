package domain

import (
	"github.com/golang/protobuf/ptypes/duration"
)

type CellBase struct {
	CellId  int64  `json:"cellId"`
	Version int64  `json:"version"`
	Branch  string `json:"Branch"`
}

type AddVersionReq struct {
	CellBase
	LockKey string `json:"LockKey"`
	Comment string `json:"comment"`
}

// 添加参数校验
type LockReq struct {
	CellId  int64  `json:"cellId"`
	Branch  string `json:"branch"`
	LockKey string `json:"lockKey"`
	// 这里是为post请求中duration能直接解析，所以使用了google的duration
	LockDuration duration.Duration `json:"lockDuration"`
}

/*
示例而已，因此字段只有几个
*/
type User struct {
	UserId int64
	AddVersionReq
	//etc ..
}
