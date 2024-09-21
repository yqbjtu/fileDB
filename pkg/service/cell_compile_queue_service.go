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

// WaitingToCompileQueueSize 查询等待编译的总队列长度
func (s *CellCompileQueueService) WaitingToCompileQueueSize() int64 {
	return s.bizStore.FindAllSize()
}

// WaitingToCompileQueueSizeByBranch 按照branch查询等待编译的队列长度
func (s *CellCompileQueueService) WaitingToCompileQueueSizeByBranch(branch string) int64 {
	return s.bizStore.FindAllByBranch(branch)
}
