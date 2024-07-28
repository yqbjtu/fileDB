package domain

import (
	"github.com/golang/protobuf/ptypes/duration"
)

type CellBase struct {
	CellId  int64  `json:"cellId" validate:"required,gte=1"  example:"12345"`
	Version int64  `json:"version" validate:"required,gte=1"  example:"2"`
	Branch  string `json:"Branch" validate:"required,gt=0,lte=100" binding:"required,gt=0,lte=100"`
}

type AddVersionReq struct {
	CellBase
	LockKey string `json:"LockKey" validate:"required,gt=0,lte=100" binding:"required,gt=0,lte=100"`
	Comment string `json:"comment"`
}

// 添加参数校验
type LockReq struct {
	CellId  int64  `json:"cellId"  validate:"required,gte=1"  example:"12345"`
	Branch  string `json:"branch"  validate:"required,gt=0,lte=100" binding:"required,gt=0,lte=100"`
	LockKey string `json:"lockKey" validate:"required,gt=0,lte=100" binding:"required,gt=0,lte=100"`
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
