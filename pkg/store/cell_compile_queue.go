package store

import (
	"fileDB/pkg/domain"
	"fileDB/pkg/log"
	"fmt"
	"gorm.io/gorm"
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

func (s *CellCompileQueueStore) Delete(cellId, version int64, branch string) (domain.CellCompileQueue, error) {
	var objResult domain.CellCompileQueue
	result := s.db.Delete(&objResult, "cell_id = ? and version =? and branch = ?", cellId, version, branch)
	if result.Error != nil {
		return objResult, fmt.Errorf("delete specific CellCompileQueue failed, err:%v", result.Error)
	}

	return objResult, nil
}

func (s *CellCompileQueueStore) FindAllSize() int64 {
	var count int64
	s.db.Model(&domain.CellCompileQueue{}).Count(&count)
	return count
}

func (s *CellCompileQueueStore) FindAllByBranch(branch string) int64 {
	var count int64
	s.db.Model(&domain.CellCompileQueue{}).Where("branch = ?", branch).Count(&count)
	return count
}

// Upsert if the cellId+branch already exists, update the record, otherwise insert a new record
func (s *CellCompileQueueStore) Upsert(obj domain.CellCompileQueue) (*domain.CellCompileQueue, error) {

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
			log.Errorf("failed to update cell compile queue, err:%v", result.Error)
			return nil, fmt.Errorf("fail to cell compile queue, err:%v", result.Error)
		}

		return &objResult, nil
	}

	// insert a new record
	result = s.db.Save(&obj)
	if result.Error != nil {
		log.Errorf("failed to insert cell compile queue, err:%v", result.Error)
		return nil, fmt.Errorf("fail to cell compile queue, err:%v", result.Error)
	}

	return nil, nil
}
