package store

import (
	"fileDB/pkg/domain"
	"fmt"
	"gorm.io/gorm"
	"k8s.io/klog"
)

type CellCompileQueueStore struct {
	db *gorm.DB
}

func NewCellCompileQueueStore(db *gorm.DB) *CellCompileQueueStore {
	return &CellCompileQueueStore{
		db: db,
	}
}

func (s *CellCompileQueueStore) Find(cellId int64, branch string) (domain.CellCompileQueue, error) {
	var objResult domain.CellCompileQueue
	result := s.db.Find(&objResult, "cell_id = ? and branch = ?", cellId, branch)
	if result.Error != nil {
		return objResult, fmt.Errorf("find cell history failed, err:%v", result.Error)
	}

	return objResult, nil
}

// Insert if the cellId+branch already exists, update the record, otherwise insert a new record
func (s *CellCompileQueueStore) Insert(obj domain.CellCompileQueue) (*domain.CellCompileQueue, error) {

	// find the record by branch cellId , if it exists, update the record, otherwise insert a new record
	var objResult domain.CellCompileQueue
	result := s.db.Find(&objResult, "cell_id = ? and branch = ?", obj.CellId, obj.Branch)
	if result.Error != nil {
		return nil, fmt.Errorf("find cell compile queue failed, err:%v", result.Error)
	}

	if objResult.ID > 0 {
		// update the record
		result := s.db.Model(&objResult).Updates(&obj)
		if result.Error != nil {
			klog.Errorf("failed to update cell compile queue, err:%v", result.Error)
			return nil, fmt.Errorf("fail to cell compile queue, err:%v", result.Error)
		}

		return &objResult, nil
	}

	// insert a new record
	result = s.db.Save(&obj)
	if result.Error != nil {
		klog.Errorf("failed to insert cell compile queue, err:%v", result.Error)
		return nil, fmt.Errorf("fail to cell compile queue, err:%v", result.Error)
	}

	return nil, nil
}
