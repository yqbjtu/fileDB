package store

import (
	"fileDB/pkg/domain"
	"fmt"
	"gorm.io/gorm"
	"k8s.io/klog"
)

type CellGisMetaStore struct {
	db *gorm.DB
}

func NewCellGisMetaStore(db *gorm.DB) *CellGisMetaStore {
	return &CellGisMetaStore{
		db: db,
	}
}

func (s *CellGisMetaStore) Find(cellId int64, branch string) (domain.CellGisMeta, error) {
	var cellGisMeta domain.CellGisMeta
	result := s.db.Find(&cellGisMeta, "cell_id = ? and branch = ?", cellId, branch)
	if result.Error != nil {
		return cellGisMeta, fmt.Errorf("find cell history failed, err:%v", result.Error)
	}

	return cellGisMeta, nil
}

func (s *CellGisMetaStore) Insert(cellGisMeta domain.CellGisMeta) (*domain.CellGisMeta, error) {

	result := s.db.Save(&cellGisMeta)
	if result.Error != nil {
		klog.Errorf("failed to insert cell gis meta, err:%v", result.Error)
		return nil, fmt.Errorf("fail to save cell history, err:%v", result.Error)
	}

	return nil, nil
}
