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

func (s *CellCompileQueueService) Find(cellId int64, branch string) (domain.CellCompileQueue, error) {
	var cellInQueue domain.CellCompileQueue
	cellInQueue, err := s.bizStore.Find(cellId, branch)
	if err != nil {
		return cellInQueue, err
	}
	return cellInQueue, nil
}

// Upsert 每个branch+cellId在队列中只有一个元素， 如有新的版本产生，会覆盖之前的版本
func (s *CellCompileQueueService) Upsert(cellCompileQueue domain.CellCompileQueue) error {
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
