package domain

type AddVersionReq struct {
	CellId    string
	Version   int64
	Namespace string
	LockKey   string
	Comment   string
}

/*
示例而已，因此字段只有几个
*/
type User struct {
	UserId int64
	AddVersionReq
	//etc ..
}
