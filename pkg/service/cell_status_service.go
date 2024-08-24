package service

import (
	"fileDB/pkg/domain"
	"fileDB/pkg/store"
)

type CellStatusService struct {
	bizStore *store.CellStatusStore
}

func NewCellStatusService(bizStore *store.CellStatusStore) *CellStatusService {
	return &CellStatusService{
		bizStore: bizStore,
	}
}

func (s *CellStatusService) Find(cellId int64, branch string) (domain.CellStatus, error) {
	return s.bizStore.Find(cellId, branch)
}

func (s *CellStatusService) Insert(cellStatus domain.CellStatus) (*domain.CellStatus, error) {
	return s.bizStore.Save(cellStatus)
}
