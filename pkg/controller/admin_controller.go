package controller

import (
	"fileDB/pkg/config"
	mydomain "fileDB/pkg/domain"
	"fileDB/pkg/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AdminController struct {
	// service or some to access DB method
	globalConfig        *config.GlobalConfig
	cellCvsSvc          *service.CellCvsService
	cellHistorySvc      *service.CellHistoryService
	cellCompileQueueSvc *service.CellCompileQueueService
}

func NewAdminController(globalConfig *config.GlobalConfig,
	cellCvsSvc *service.CellCvsService,
	cellHistorySvc *service.CellHistoryService,
	cellCompileQueueSvc *service.CellCompileQueueService) *AdminController {
	controller := AdminController{
		globalConfig:        globalConfig,
		cellCvsSvc:          cellCvsSvc,
		cellHistorySvc:      cellHistorySvc,
		cellCompileQueueSvc: cellCompileQueueSvc,
	}
	return &controller
}

// FindAllWaitingToCompileQueue 查询总的等待编编译的队列长度
// @Summary FindAllWaitingToCompileQueue 查询总的等待编编译的队列长度
// @Description query all waiting to compile queue
// @Tags query
// @Accept  json
// @Produce json
// @Success 200 {object} mydomain.CommonResult "ok"
// @Router /api/v1/admin/compileQueueSize [get]
func (c *AdminController) FindAllWaitingToCompileQueue(ctx *gin.Context) {
	count := c.cellCompileQueueSvc.WaitingToCompileQueueSize()
	CommonResult := mydomain.CommonResult{Code: 0, Data: count, Msg: "done"}
	ctx.JSON(http.StatusOK, CommonResult)
	return
}

func (c *AdminController) FindAllWaitingToCompileQueueByBranch(ctx *gin.Context) {
	branchStr := ctx.Query("branch")
	count := c.cellCompileQueueSvc.WaitingToCompileQueueSizeByBranch(branchStr)
	msg := fmt.Sprintf("branch '%s' has %d cell to compile", branchStr, count)
	CommonResult := mydomain.CommonResult{Code: 0, Data: count, Msg: msg}
	ctx.JSON(http.StatusOK, CommonResult)
	return
}
