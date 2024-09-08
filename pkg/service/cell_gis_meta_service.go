package service

import (
	"fileDB/pkg/config"
	"fileDB/pkg/domain"
	"fileDB/pkg/log"
	"fileDB/pkg/store"
	"fmt"
)

type CellGisMetaService struct {
	globalConfig        *config.GlobalConfig
	cellHistorySvc      *CellHistoryService
	cellCompileQueueSvc *CellCompileQueueService
	cellStatusSvc       *CellStatusService
	bizStore            *store.CellGisMetaStore
}

func NewCellGisMetaService(
	globalConfig *config.GlobalConfig,
	cellHistorySvc *CellHistoryService,
	cellCompileQueueSvc *CellCompileQueueService,
	cellStatusSvc *CellStatusService,
	bizStore *store.CellGisMetaStore) *CellGisMetaService {
	return &CellGisMetaService{
		globalConfig:        globalConfig,
		cellHistorySvc:      cellHistorySvc,
		cellCompileQueueSvc: cellCompileQueueSvc,
		cellStatusSvc:       cellStatusSvc,
		bizStore:            bizStore,
	}
}

func (s *CellGisMetaService) UpsertGisMeta(req domain.CellGisMeta) (domain.CommonResult, error) {
	cellMeta, err := s.bizStore.Find(req.CellId, req.Branch)
	if err != nil {
		log.Errorf("failed to find cell status, err:%v", err)
	}

	// if cell not exist, create a new cell status
	if cellMeta.CellId == 0 {
		_, err = s.bizStore.Insert(req)
		if err != nil {
			log.Errorf("failed to save cell status, err:%v", err)
			CommonResult := domain.CommonResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to  save cell status, err:%v", err)}
			return CommonResult, nil
		} else {
			CommonResult := domain.CommonResult{Code: 0, Data: req, Msg: "add the first version ok"}
			return CommonResult, nil
		}
	}

	CommonResult := domain.CommonResult{Code: 0, Data: nil, Msg: fmt.Sprintf("cell %d add new version done", req.CellId)}
	return CommonResult, nil
}

func (s *CellGisMetaService) UnLock(req domain.AddVersionReq) (domain.CommonResult, error) {

	return domain.CommonResult{}, nil
}
