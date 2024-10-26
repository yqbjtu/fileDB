package controller

import (
	"fileDB/pkg/domain"
	"fileDB/pkg/service"
	"fileDB/pkg/store"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/klog"
	"net/http"
)

type CompileQueueController struct {
	// service or some to access DB method
	cellHistorySvc      *service.CellHistoryService
	cellStatusStore     *store.CellStatusStore
	cellCompileQueueSvc *service.CellCompileQueueService
}

func NewCompileQueueController(cellHistorySvc *service.CellHistoryService,
	cellStatusStore *store.CellStatusStore,
	cellCompileQueueSvc *service.CellCompileQueueService) *CompileQueueController {
	controller := CompileQueueController{
		cellHistorySvc:      cellHistorySvc,
		cellStatusStore:     cellStatusStore,
		cellCompileQueueSvc: cellCompileQueueSvc,
	}
	return &controller
}

// WaitingToCompileQueueSize @Summary get  queue size of cell to be compiled
func (c *CompileQueueController) WaitingToCompileQueueSize(ctx *gin.Context) {
	size := c.cellCompileQueueSvc.WaitingToCompileQueueSize()

	commonRes :=
		domain.CommonResult{Code: 1, Data: size, Msg: "total queue size to be compiled for all branches"}

	ctx.JSON(http.StatusOK, commonRes)
	return
}

// FindAllWaitingToCompileQueue 查询总的等待编编译的队列长度
// @Summary FindAllWaitingToCompileQueue 查询总的等待编编译的队列长度
// @Description query all waiting to compile queue
// @Tags query
// @Accept  json
// @Produce json
// @Success 200 {object} mydomain.CommonResult "ok"
// @Router /api/v1/admin/compileQueueSize [get]
func (c *CompileQueueController) WaitingToCompileQueueSizeByBranch(ctx *gin.Context) {
	branchStr := ctx.Query("branch")
	if branchStr == "" {
		klog.Errorf("branch '%s' can't be empty", branchStr)
		msg := fmt.Sprintf("query paramenter branch can't be empty")
		ctx.JSON(http.StatusBadRequest, domain.NewErrorRespWithMsg(-1, msg))
	}

	size := c.cellCompileQueueSvc.WaitingToCompileQueueSizeByBranch(branchStr)
	msg := fmt.Sprintf("branch '%s' has %d cell to compile", branchStr, size)
	commonRes :=
		domain.CommonResult{Code: 1, Data: size, Msg: msg}

	ctx.JSON(http.StatusOK, commonRes)
	return
}
