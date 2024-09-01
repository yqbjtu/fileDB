package controller

import (
	mydomain "fileDB/pkg/domain"
	"fileDB/pkg/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

type MiscController struct {
	// service or some to access DB method
}

func NewMiscController() *MiscController {
	controller := MiscController{}
	return &controller
}

func (c *MiscController) FreeMemory(context *gin.Context) {
	debug.FreeOSMemory()

	CommonResult := mydomain.CommonResult{Code: 0, Data: nil, Msg: "free mem ok"}
	context.JSON(http.StatusOK, CommonResult)
}

func (c *MiscController) BuildInfo(context *gin.Context) {
	log.Infof("build info")

	CommonResult := mydomain.CommonResult{Code: 0, Data: "v1.0", Msg: "ok"}
	context.JSON(http.StatusOK, CommonResult)
}
