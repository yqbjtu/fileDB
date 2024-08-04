package domain

type CellGisMeta struct {
	BaseModel
	CellId  int64
	Branch  string
	Version int64
	MinX    float64
	MinY    float64
	MaxX    float64
	MaxY    float64
}

func (CellGisMeta) TableName() string {
	return "cell_gis_meta"
}
