package controller

import (
	"fileDB/pkg/config"
	mydomain "fileDB/pkg/domain"
	"fileDB/pkg/service"
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

// Backup @Summary backup specific branch cell files
func (c *AdminController) Backup(ctx *gin.Context) {
	// only start the branch backup goroutine, return immediately
	CommonResult := mydomain.CommonResult{Code: 0, Data: nil, Msg: "start to backup branch files"}
	ctx.JSON(http.StatusOK, CommonResult)
	return
}
