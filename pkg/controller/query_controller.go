package controller

import (
	mydomain "fileDB/pkg/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/klog"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
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

func (c *QueryController) DownloadFile(ctx *gin.Context) {
	var req mydomain.CellBase
	var err error
	cellIdStr := ctx.Query("cellId")
	versionStr := ctx.Query("version")
	namespaceStr := ctx.Query("namespace")

	if cellIdStr == "" || namespaceStr == "" || versionStr == "" {
		klog.Errorf("cellId '%s' can't be empty", cellIdStr)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errMsg": "cellId/version/ is empty",
		})
		return
	} else {
		req.Version, err = strconv.ParseInt(versionStr, 10, 32)
		if err != nil {
			klog.Errorf("failed to convert (%s)to int64, err:%v", versionStr, err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errMsg": "version is not int",
			})
			return
		}

		req.CellId = cellIdStr
		req.Namespace = namespaceStr
	}

	// 你可以访问header来获取文件名称、文件大小和文件类型等信息
	filename := fmt.Sprintf("%s@@%s@@%d.osm", req.CellId, req.Namespace, req.Version)
	// 定义文件保存路径
	cellPath := fmt.Sprintf("/tmp/osmdb/data/%s/", req.Namespace) + filename

	// 先检测该cellPath是否存在，如果不存在报错
	// 如果目录不存在，则创建改目录
	if _, err := os.Stat(cellPath); err != nil {
		commentResult :=
			mydomain.CommentResult{Code: -1, Data: nil,
				Msg: fmt.Sprintf("fail to find file %s, err:%v", filename, err)}
		if os.IsNotExist(err) {
			ctx.JSON(http.StatusNotFound, commentResult)
			return
		} else {
			// 如果检查时发生其他错误，则返回错误信息
			ctx.JSON(http.StatusInternalServerError, commentResult)
			return
		}
	}

	ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
	// 发送文件内容给客户端
	ctx.File(cellPath)
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
