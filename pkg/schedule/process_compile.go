package schedule

import (
	"fileDB/pkg/service"
	"time"
)

type CellCompileScheduleService struct {
	cellCompileQueueSvc service.CellCompileQueueService
}

func NewCellCompileScheduleService(cellCompileQueueSvc service.CellCompileQueueService) *CellCompileScheduleService {
	return &CellCompileScheduleService{
		cellCompileQueueSvc: cellCompileQueueSvc,
	}
}

func (s *CellCompileScheduleService) ProcessCompile(cellCompileQueueSvc service.CellCompileQueueService) {
	ticker := time.NewTicker(3 * time.Second)

	// 使用for循环和select语句来实现定时任务
	for {
		select {
		case <-ticker.C:
			s.ProcessToCompileCellQueue()
		}
	}
}

func (s *CellCompileScheduleService) ProcessToCompileCellQueue() {
	// 处理内容
}
