package controller

import (
	"fileDB/pkg/config"
	mydomain "fileDB/pkg/domain"
	"fileDB/pkg/log"
	"fileDB/pkg/service"
	"fileDB/pkg/store"
	"fileDB/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type QueryController struct {
	// service or some to access DB method
	cellHistorySvc  *service.CellHistoryService
	cellStatusStore *store.CellStatusStore
}

func NewQueryController(cellHistorySvc *service.CellHistoryService, cellStatusStore *store.CellStatusStore) *QueryController {
	controller := QueryController{
		cellHistorySvc:  cellHistorySvc,
		cellStatusStore: cellStatusStore,
	}
	return &controller
}

// @Summary query cell status
// @Description check cell exist or not, cell is checkout or not etc
// @Tags query
// @Accept  json
// @Produce json
// @Success 200 {string} string	"ok"
// @Failure 400 {string} string "We need cellId and branch"
// @Router /api/v1/cellversion/status [get]
func (c *QueryController) CellStatus(ctx *gin.Context) {
	req, err := util.GetCellBaseFromParameter(ctx, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errMsg": err.Error(),
		})
		return
	}

	cellStatus, err := c.cellStatusStore.Find(req.CellId, req.Branch)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"errMsg": err.Error(),
		})
		return
	}

	if cellStatus.CellId == 0 {
		log.Infof("cell not exist")
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

	req, err = util.GetCellBaseFromParameter(ctx, true)
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
	log.Infof("build info")
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
	req, err := util.GetCellBaseFromParameter(ctx, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errMsg": err.Error(),
		})
		return
	}

	var commentResult mydomain.CommentResult
	historyList, err := c.cellHistorySvc.Find(req.CellId, req.Branch)
	if err != nil {
		commentResult =
			mydomain.CommentResult{Code: -1, Data: nil,
				Msg: fmt.Sprintf("fail to find history branch:%s, id:%d, err:%v", req.Branch, req.CellId, err)}
		ctx.JSON(http.StatusInternalServerError, commentResult)
	} else {
		commentResult =
			mydomain.CommentResult{Code: -1, Data: nil,
				Msg: fmt.Sprintf("branch:%s, cellId:%d, size:%d", req.Branch, req.CellId, len(historyList))}
		ctx.JSON(http.StatusOK, commentResult)
	}

	return
}
