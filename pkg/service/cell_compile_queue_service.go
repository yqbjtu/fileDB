package service

import (
	"fileDB/pkg/domain"
	"fileDB/pkg/store"
)

type CellCompileQueueService struct {
	bizStore *store.CellCompileQueueStore
}

func NewCellCompileQueueService(bizStore *store.CellCompileQueueStore) *CellCompileQueueService {
	return &CellCompileQueueService{
		bizStore: bizStore,
	}
}

func (s *CellCompileQueueService) Find(cellId int64, branch string) ([]domain.CellCompileQueue, error) {

	var cellInQueueList []domain.CellCompileQueue

	return cellInQueueList, nil
}

func (s *CellCompileQueueService) Insert(cellCompileQueue domain.CellCompileQueue) error {
	_, err := s.bizStore.Upsert(cellCompileQueue)
	if err != nil {
		return err
	}

	return nil
}
