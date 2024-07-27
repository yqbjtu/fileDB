package store

import (
	"fileDB/pkg/domain"
	"fmt"
	"gorm.io/gorm"
)

type CellStatusStore struct {
	db *gorm.DB
}

func NewCellStatusStore(db *gorm.DB) *CellStatusStore {
	return &CellStatusStore{
		db: db,
	}
}

func (s *CellStatusStore) Find(cellId int64, branch string) (domain.CellStatus, error) {

	var cellStatus domain.CellStatus
	result := s.db.Find(&cellStatus, "cell_id = ? and branch = ?", cellId, branch)
	if result.Error != nil {
		return cellStatus, fmt.Errorf("find cell status failed, err:%v", result.Error)
	}

	return cellStatus, nil
}
