package controller

import (
	"fileDB/pkg/config"
	mydomain "fileDB/pkg/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/klog"
	"net/http"
	"os"
	"strconv"
)

type QueryController struct {
	// service or some to access DB method
}

func NewQueryController() *QueryController {
	controller := QueryController{}
	return &controller
}

func getCellBaseFromParameter(ctx *gin.Context) (mydomain.CellBase, error) {
	var req mydomain.CellBase
	//var err error
	cellIdStr := ctx.Query("cellId")
	branchStr := ctx.Query("branch")
	if cellIdStr == "" || branchStr == "" {
		klog.Errorf("cellId '%s'/branch '%s' can't be empty", cellIdStr, branchStr)
		return req, fmt.Errorf("cellId or branch is empty")
	} else {
		cellId, err := strconv.ParseInt(cellIdStr, 10, 64)
		if err != nil {
			klog.Errorf("failed to convert (%s)to int64, err:%v", cellIdStr, err)
			return req, fmt.Errorf("cellId is int type")
		}
		req.CellId = cellId
		req.Branch = branchStr
	}

	return req, nil
}

// @Summary query cell status
// @Description check cell exist or not, cell is checkout or not etc
// @Tags query
// @Accept  json
// @Produce json
// @Success 200 {string} string	"ok"
// @Failure 400 {string} string "We need cellId and branch"
// @Router /api/v1/cellversion/status [get]
func (c *QueryController) FileStatus(ctx *gin.Context) {
	req, err := getCellBaseFromParameter(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errMsg": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  fmt.Sprintf("cellId:%d, branch:%s", req.CellId, req.Branch),
	})
}

// @Summary download specific version cell file
// @Description download cell file by cellId, version and branch
// @Tags query
// @Accept  json
// @Produce json
// @Success 200 {string} string	"ok"
// @Failure 400 {string} string "We need cellId,version and branch"
// @Router /example/download [get]
func (c *QueryController) DownloadFile(ctx *gin.Context) {
	var req mydomain.CellBase
	var err error
	cellIdStr := ctx.Query("cellId")
	versionStr := ctx.Query("version")
	branchStr := ctx.Query("branch")

	if cellIdStr == "" || branchStr == "" || versionStr == "" {
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

		cellId, err := strconv.ParseInt(cellIdStr, 10, 64)
		if err != nil {
			klog.Errorf("failed to convert (%s)to int64, err:%v", cellIdStr, err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errMsg": "cellId is int type",
			})
			return
		}
		req.CellId = cellId
		req.Branch = branchStr
	}

	// 你可以访问header来获取文件名称、文件大小和文件类型等信息
	filename := fmt.Sprintf("%d@@%s@@%d.osm", req.CellId, req.Branch, req.Version)
	// 定义文件保存路径
	baseOsmDataDir := config.GetConfig().OSMConfig.DataDir
	cellPath := fmt.Sprintf("%s/%s/", baseOsmDataDir, req.Branch) + filename

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

func (c *QueryController) FileBBoxInfo(ctx *gin.Context) {
	klog.Infof("build info")
	//H is a shortcut for map[string]interface{}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": nil,
		"msg":  "ok",
	})
}

// query the cell history , such as addVersion,
func (c *QueryController) History(ctx *gin.Context) {
	var req mydomain.CellBase
	req, err := getCellBaseFromParameter(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errMsg": err.Error(),
		})
		return
	}

	// query db to find the history
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": nil,
		"msg":  fmt.Sprintf("cellId:%d, branch:%s", req.CellId, req.Branch),
	})
}
