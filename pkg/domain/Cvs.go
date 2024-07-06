package domain

type CellBase struct {
	CellId    string
	Version   int64
	Namespace string
}

type AddVersionReq struct {
	CellBase
	LockKey string
	Comment string
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
