package service

import (
	"fileDB/pkg/domain"
	"fileDB/pkg/store"
)

type CellHistoryService struct {
	bizStore *store.CellHistoryStore
}

func NewCellHistoryService(bizStore *store.CellHistoryStore) *CellHistoryService {
	return &CellHistoryService{
		bizStore: bizStore,
	}
}

func (s *CellHistoryService) Find(cellId int64, branch string) ([]domain.CellHistory, error) {

	var cellHistoryList []domain.CellHistory

	return cellHistoryList, nil
}

func (s *CellHistoryService) Insert(history domain.CellHistory) error {
	_, err := s.bizStore.Insert(history)
	if err != nil {
		return err
	}

	return nil
}
