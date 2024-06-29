package controller

import (
	"github.com/gin-gonic/gin"
	"k8s.io/klog"
	"net/http"
	"runtime/debug"
)

type QueryController struct {
	// service or some to access DB method
}

func NewQueryController() *QueryController {
	controller := QueryController{}
	return &controller
}

func (c *QueryController) FileStatus(context *gin.Context) {
	debug.FreeOSMemory()
	context.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "create user successfully",
	})
}

func (c *QueryController) DownloadFile(context *gin.Context) {
	klog.Infof("build info")
	//H is a shortcut for map[string]interface{}

	context.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": nil,
		"msg":  "ok",
	})
}

func (c *QueryController) FileBBoxInfo(context *gin.Context) {
	klog.Infof("build info")
	//H is a shortcut for map[string]interface{}

	context.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": nil,
		"msg":  "ok",
	})
}
