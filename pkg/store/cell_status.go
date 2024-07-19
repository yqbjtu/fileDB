package store

import (
	"context"
	"fileDB/pkg/domain"
	"fmt"
	"gorm.io/gorm"
)

type CellStatusStore struct {
	db *gorm.DB
}

// NewForceMatchStore -
func NewCellStatusStore(db *gorm.DB) *CellStatusStore {
	return &CellStatusStore{
		db: db,
	}
}

// Save 批量保存
func (s *CellStatusStore) Find(ctx context.Context, cellId, branch string) error {

	var cellStatus domain.CellStatus
	result := s.db.Find(&cellStatus, "cell_id = ? and branch = ?", cellId, branch)
	if result.Error != nil {
		return fmt.Errorf("find cell status failed, err:%v", result.Error)
	}

	return nil

}
