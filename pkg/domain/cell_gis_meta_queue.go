package domain

// 所有cell在生成新版本的时候在该表插入一条记录，表示需要编译该版本。如果该表中对应branch和cellId已经存在过，直接更新老的， 编译完毕后直接删除该记录
type CellGisMetaQueue struct {
	BaseModel
	CellId  int64
	Branch  string
	Version int64
	// 如果一个cell的上一个版本没有编译，直接将Version更新为最新需要编译的，将上一个版本记录到LastVersion中
	LastVersion int64
	// 开始的时候优先级都一样，按照时间顺序，如果临时某个branch需要提升优先级，将该branch的priority升级到1
	Priority int64
}

func (CellGisMetaQueue) TableName() string {
	return "cell_gis_meta_queue"
}
