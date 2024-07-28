package store

import (
	"fileDB/pkg/domain"
	"fmt"
	"gorm.io/gorm"
)

type CellHistoryStore struct {
	db *gorm.DB
}

func NewCellHistoryStore(db *gorm.DB) *CellHistoryStore {
	return &CellHistoryStore{
		db: db,
	}
}

func (s *CellHistoryStore) Find(cellId int64, branch string) ([]domain.CellHistory, error) {

	var cellHistoryList []domain.CellHistory
	result := s.db.Find(&cellHistoryList, "cell_id = ? and branch = ?", cellId, branch)
	if result.Error != nil {
		return cellHistoryList, fmt.Errorf("find cell history failed, err:%v", result.Error)
	}

	return cellHistoryList, nil
}
