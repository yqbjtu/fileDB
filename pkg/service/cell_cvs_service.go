package service

import (
	"fileDB/pkg/common"
	"fileDB/pkg/config"
	"fileDB/pkg/domain"
	"fileDB/pkg/log"
	"fmt"
	"time"
)

type CellCvsService struct {
	globalConfig        *config.GlobalConfig
	cellHistorySvc      *CellHistoryService
	cellCompileQueueSvc *CellCompileQueueService
	cellStatusSvc       *CellStatusService
}

func NewCellCvsService(
	globalConfig *config.GlobalConfig,
	cellHistorySvc *CellHistoryService,
	cellCompileQueueSvc *CellCompileQueueService,
	cellStatusSvc *CellStatusService) *CellCvsService {
	return &CellCvsService{
		globalConfig:        globalConfig,
		cellHistorySvc:      cellHistorySvc,
		cellCompileQueueSvc: cellCompileQueueSvc,
		cellStatusSvc:       cellStatusSvc,
	}
}

func (s *CellCvsService) AddNewVersion(req domain.AddVersionReq) (domain.CommonResult, error) {
	cellStatus, err := s.cellStatusSvc.Find(req.CellId, req.Branch)
	if err != nil {
		log.Errorf("failed to find cell status, err:%v", err)
	}

	// if cell not exist, create a new cell status
	if cellStatus.CellId == 0 {
		cellStatus.CellId = req.CellId
		cellStatus.LatestVersion = req.Version
		cellStatus.LockKey = ""
		cellStatus.Branch = req.Branch
		_, err = s.cellStatusSvc.Insert(cellStatus)
		if err != nil {
			log.Errorf("failed to save cell status, err:%v", err)
			commonRes := domain.CommonResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to  save cell status, err:%v", err)}
			return commonRes, nil
		} else {
			commonRes := domain.CommonResult{Code: 0, Data: req, Msg: "add the first version ok"}
			return commonRes, nil
		}
	}

	// the req.Version should be the latest version + 1
	expectedVersion := cellStatus.LatestVersion + 1
	if req.Version != expectedVersion {
		errMsg := fmt.Sprintf("cellId:%d, current latest version is %d, expectedVersion should be %d, not %d", req.CellId,
			cellStatus.LatestVersion, expectedVersion, req.Version)
		commonRes := domain.CommonResult{Code: -1, Data: nil, Msg: errMsg}
		return commonRes, fmt.Errorf(errMsg)
	}

	// the cell should not be locked, or it is locked by req.LockKey
	if cellStatus.LockKey != "" && cellStatus.LockKey != req.LockKey {
		errMsg := fmt.Sprintf("cellId:%d is locked by %q, not %q", req.CellId, cellStatus.LockKey, req.LockKey)
		commonRes := domain.CommonResult{Code: -1, Data: nil, Msg: errMsg}
		return commonRes, fmt.Errorf(errMsg)
	}

	// update the cell status with latestVersion
	cellStatus.LatestVersion = req.Version
	cellStatus.LockKey = ""
	_, err = s.cellStatusSvc.Insert(cellStatus)
	if err != nil {
		log.Errorf("failed to save cell status, err:%v", err)
		commonRes := domain.CommonResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to SaveUploadedFile, err:%v", err)}
		return commonRes, nil
	}

	cellHistory := domain.CellHistory{
		CellId:      req.CellId,
		Branch:      req.Branch,
		Version:     req.Version,
		RequestType: "CheckinRequest",
		LockKey:     req.LockKey,
		Who:         "tester1",
	}

	err = s.cellHistorySvc.Insert(cellHistory)
	if err != nil {
		commonRes := domain.CommonResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to SaveHistoryRecord, err:%v", err)}
		return commonRes, nil
	}

	// insert a event to cell_compile_queue
	cellCompileQueue := domain.CellCompileQueue{
		CellId:   req.CellId,
		Branch:   req.Branch,
		Version:  req.Version,
		Priority: 0,
	}

	err = s.cellCompileQueueSvc.Upsert(cellCompileQueue)
	if err != nil {
		log.Errorf("failed to save cell compile queue, err:%v", err)
		commonRes :=
			domain.CommonResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to save cell compile queue, err:%v", err)}
		return commonRes, common.ErrDBOperationFailure
	}

	commonRes := domain.CommonResult{Code: 0, Data: nil, Msg: fmt.Sprintf("cell %d add new version done", req.CellId)}
	return commonRes, nil
}

func (s *CellCvsService) Lock(lockReq *domain.LockReq) (domain.CommonResult, error) {
	cellStatus, err := s.cellStatusSvc.Find(lockReq.CellId, lockReq.Branch)
	if err != nil {
		log.Errorf("failed to find cell status, err:%v", err)
	}

	if cellStatus.LockKey != "" && cellStatus.LockKey != lockReq.LockKey {
		errMsg := fmt.Sprintf("cell %d has already been locked by %s now, so it can't be locked by %s again",
			lockReq.CellId, cellStatus.LockKey, lockReq.LockKey)
		commonRes := domain.CommonResult{Code: -1, Data: nil, Msg: errMsg}
		return commonRes, fmt.Errorf(errMsg)
	}
	if lockReq.LockDuration.GetSeconds() <= 10 {
		errMsg := fmt.Sprintf("cell lock duration should be gt 10s, but it is %v", lockReq.LockDuration)
		commonRes := domain.CommonResult{Code: -1, Data: nil, Msg: errMsg}
		return commonRes, fmt.Errorf(errMsg)
	}

	cellStatus.LockKey = lockReq.LockKey
	// cellStatus.LockTimeFrom等于当前时间
	fromTime := time.Now()
	cellStatus.LockTimeFrom = &fromTime

	// cellStatus.LockTimeTo等于当前时间加上一个小时
	goDuration := time.Duration(lockReq.LockDuration.GetSeconds())*time.Second + time.Duration(lockReq.LockDuration.GetNanos())*time.Nanosecond
	toTime := time.Now().Add(goDuration)
	cellStatus.LockTimeTo = &toTime
	if cellStatus.CellId == 0 {
		cellStatus.CellId = lockReq.CellId
		cellStatus.Branch = lockReq.Branch
		// 没有添加过，版本就为0
		cellStatus.LatestVersion = 0
	}

	_, err = s.cellStatusSvc.Insert(cellStatus)
	if err != nil {
		log.Errorf("failed to save cell status, err:%v", err)
		commonRes :=
			domain.CommonResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to save cell status lock info, err:%v", err)}
		return commonRes, common.ErrDBOperationFailure
	}

	customTimeFormat := "2006-01-02 15:04:05"
	// add lock record in db
	response := map[string]interface{}{
		"id":           lockReq.CellId,
		"latestVer":    cellStatus.LatestVersion,
		"branch":       lockReq.Branch,
		"lockKey":      lockReq.LockKey,
		"lockTimeFrom": cellStatus.LockTimeFrom.Format(customTimeFormat),
		"lockTimeTo":   cellStatus.LockTimeTo.Format(customTimeFormat),
	}

	return *domain.NewSuccessResp(response), nil
}

func (s *CellCvsService) UnLock(req domain.AddVersionReq) (domain.CommonResult, error) {

	return domain.CommonResult{}, nil
}
