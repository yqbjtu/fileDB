package controller

import (
	mydomain "fileDB/pkg/domain"
	"github.com/gin-gonic/gin"
	"k8s.io/klog"
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

	commentResult := mydomain.CommentResult{Code: 0, Data: nil, Msg: "free mem ok"}
	context.JSON(http.StatusOK, commentResult)
}

func (c *MiscController) BuildInfo(context *gin.Context) {
	klog.Infof("build info")

	commentResult := mydomain.CommentResult{Code: 0, Data: "v1.0", Msg: "ok"}
	context.JSON(http.StatusOK, commentResult)
}
