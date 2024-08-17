package store

import (
	"fileDB/pkg/domain"
	"fmt"
	"gorm.io/gorm"
	"k8s.io/klog"
)

type CellStatusStore struct {
	db *gorm.DB
}

func NewCellStatusStore(db *gorm.DB) *CellStatusStore {
	return &CellStatusStore{
		db: db,
	}
}

func (s *CellStatusStore) Save(cellStatus domain.CellStatus) (*domain.CellStatus, error) {

	result := s.db.Save(&cellStatus)
	if result.Error != nil {
		klog.Errorf("failed to save cell status, err:%v", result.Error)
		return nil, result.Error
	}

	return &cellStatus, nil
}

func (s *CellStatusStore) Find(cellId int64, branch string) (domain.CellStatus, error) {

	var cellStatus domain.CellStatus
	result := s.db.Find(&cellStatus, "cell_id = ? and branch = ?", cellId, branch)
	if result.Error != nil {
		return cellStatus, fmt.Errorf("find cell status failed, err:%v", result.Error)
	}

	return cellStatus, nil
}

func (s *CellStatusStore) FindAll() ([]domain.CellStatus, error) {
	// find all
	var cellStatus []domain.CellStatus
	result := s.db.Find(&cellStatus)
	if result.Error != nil {
		return cellStatus, fmt.Errorf("find cell status failed, err:%v", result.Error)
	}

	return cellStatus, nil
}
