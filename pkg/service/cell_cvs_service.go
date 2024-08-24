package service

import (
	"fileDB/pkg/config"
	"fileDB/pkg/domain"
	"fmt"
	"k8s.io/klog"
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

func (s *CellCvsService) AddNewVersion(req domain.AddVersionReq) (domain.CommentResult, error) {

	cellStatus, err := s.cellStatusSvc.Find(req.CellId, req.Branch)
	if err != nil {
		klog.Errorf("failed to find cell status, err:%v", err)
	}

	// if cell not exist, create a new cell status
	if cellStatus.CellId == 0 {
		cellStatus.CellId = req.CellId
		cellStatus.LatestVersion = req.Version
		cellStatus.LockKey = ""
		cellStatus.Branch = req.Branch
		_, err = s.cellStatusSvc.Insert(cellStatus)
		if err != nil {
			klog.Errorf("failed to save cell status, err:%v", err)
			commentResult := domain.CommentResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to  save cell status, err:%v", err)}
			return commentResult, nil
		} else {
			commentResult := domain.CommentResult{Code: 0, Data: req, Msg: "add the first version ok"}
			return commentResult, nil
		}
	}

	// the req.Version should be the latest version + 1
	expectedVersion := cellStatus.LatestVersion + 1
	if req.Version != expectedVersion {
		errMsg := fmt.Sprintf("cellId:%d, current latest version is %d, expectedVersion should be %d, not %d", req.CellId,
			cellStatus.LatestVersion, expectedVersion, req.Version)
		commentResult := domain.CommentResult{Code: -1, Data: nil, Msg: errMsg}
		return commentResult, fmt.Errorf(errMsg)
	}

	// the cell should not be locked, or it is locked by req.LockKey
	if cellStatus.LockKey != "" && cellStatus.LockKey != req.LockKey {
		errMsg := fmt.Sprintf("cellId:%d is locked by %q, not %q", req.CellId, cellStatus.LockKey, req.LockKey)
		commentResult := domain.CommentResult{Code: -1, Data: nil, Msg: errMsg}
		return commentResult, fmt.Errorf(errMsg)
	}

	// update the cell status with latestVersion
	cellStatus.LatestVersion = req.Version
	cellStatus.LockKey = ""
	_, err = s.cellStatusSvc.Insert(cellStatus)
	if err != nil {
		klog.Errorf("failed to save cell status, err:%v", err)
		commentResult := domain.CommentResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to SaveUploadedFile, err:%v", err)}
		return commentResult, nil
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
		commentResult := domain.CommentResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to SaveHistoryRecord, err:%v", err)}
		return commentResult, nil
	}

	commentResult := domain.CommentResult{Code: 0, Data: nil, Msg: fmt.Sprintf("cell %d add new version done", req.CellId)}
	return commentResult, nil
}

func (s *CellCvsService) Lock(history domain.CellHistory) error {

	return nil
}
