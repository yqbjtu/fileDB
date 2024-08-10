package controller

import (
	"fileDB/pkg/config"
	mydomain "fileDB/pkg/domain"
	"fileDB/pkg/store"
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

func getCellBaseFromParameter(ctx *gin.Context, isWithVersion bool) (mydomain.CellBase, error) {
	var req mydomain.CellBase
	var err error
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
			return req, fmt.Errorf("cellId should be int type")
		}
		req.CellId = cellId
		req.Branch = branchStr
	}

	if isWithVersion {
		versionStr := ctx.Query("version")
		if versionStr == "" {
			klog.Errorf("version '%s' can't be empty", versionStr)
			return req, fmt.Errorf("version is empty")
		}

		req.Version, err = strconv.ParseInt(versionStr, 10, 32)
		if err != nil {
			klog.Errorf("failed to convert (%s)to int64, err:%v", versionStr, err)
			return req, fmt.Errorf("version should be int type")
		}
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
	req, err := getCellBaseFromParameter(ctx, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errMsg": err.Error(),
		})
		return
	}

	cellStatusStore := store.NewCellStatusStore(store.MyDB)
	cellStatus, err := cellStatusStore.Find(req.CellId, req.Branch)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"errMsg": err.Error(),
		})
		return
	}

	if cellStatus.CellId == 0 {
		klog.Infof("cell not exist")
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": 0,
			"data": nil,
			"msg":  "cell not exist",
		})
		return
	}

	commentResult := mydomain.CommentResult{Code: 0, Data: cellStatus, Msg: "cell status done"}
	ctx.JSON(http.StatusOK, commentResult)
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

	req, err = getCellBaseFromParameter(ctx, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errMsg": err.Error(),
		})
		return
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

// History query the cell history , such as addVersion, lock,. unlock, etc
func (c *QueryController) History(ctx *gin.Context) {
	var req mydomain.CellBase
	req, err := getCellBaseFromParameter(ctx, false)
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
