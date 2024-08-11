package domain

// when the cell is added new version, it is added to the compile queue, where it is compiled and then delete from this queue
// one cell only has one item in the queue, the higher the priority, the earlier it is compiled.
// the higher version will override the lower version for the same cell
type CellCompileQueue struct {
	BaseModel
	CellId   int64
	Branch   string
	Version  int64
	Priority int64
}

func (CellCompileQueue) TableName() string {
	return "cell_compile_meta"
}
