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
	cellHistorySvc          *service.CellHistoryService
	cellStatusStore         *store.CellStatusStore
	CellCompileQueueService *service.CellGisMetaService
}

func NewQueryController(cellHistorySvc *service.CellHistoryService,
	cellStatusStore *store.CellStatusStore,
	cellGisMetaSvc *service.CellGisMetaService) *QueryController {
	controller := QueryController{
		cellHistorySvc:  cellHistorySvc,
		cellStatusStore: cellStatusStore,
		cellGisMetaSvc:  cellGisMetaSvc,
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
		msg := fmt.Sprintf("fail to parse parameter from url query, err:%v", err)
		ctx.JSON(http.StatusBadRequest, mydomain.NewErrorRespWithMsg(-1, msg))
		return
	}

	cellStatus, err := c.cellStatusStore.Find(req.CellId, req.Branch)
	if err != nil {
		msg := fmt.Sprintf("fail to query db, err:%v", err)
		ctx.JSON(http.StatusInternalServerError, mydomain.NewErrorRespWithMsg(-1, msg))
		return
	}

	if cellStatus.CellId == 0 {
		log.Infof("cell not exist")
		msg := fmt.Sprintf("cellId %d does not exist", req.CellId)
		ctx.JSON(http.StatusNotFound, mydomain.NewSuccessRespWithMsg(nil, msg))
		return
	}

	ctx.JSON(http.StatusOK, mydomain.NewSuccessResp(cellStatus))
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
		msg := fmt.Sprintf("fail to parse parameter from url query, err:%v", err)
		ctx.JSON(http.StatusBadRequest, mydomain.NewErrorRespWithMsg(-1, msg))
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
		msg := fmt.Sprintf("fail to find file %s, err:%v", filename, err)
		if os.IsNotExist(err) {
			ctx.JSON(http.StatusNotFound, mydomain.NewErrorRespWithMsg(-1, msg))
			return
		} else {
			// 如果检查时发生其他错误，则返回错误信息
			ctx.JSON(http.StatusInternalServerError, mydomain.NewErrorRespWithMsg(-1, msg))
			return
		}
	}

	ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
	// 发送文件内容给客户端
	ctx.File(cellPath)
}

func (c *QueryController) BBoxInfo(ctx *gin.Context) {
	req, err := util.GetCellBaseFromParameter(ctx, false)
	if err != nil {
		msg := fmt.Sprintf("fail to parse parameter from url query, err:%v", err)
		ctx.JSON(http.StatusBadRequest, mydomain.NewErrorRespWithMsg(-1, msg))
		return
	}

	commonRes := c.cellGisMetaSvc.BBoxInfo(req.Branch, req.CellId)
	ctx.JSON(http.StatusOK, commonRes)
	return
}

// History query the cell history , such as addVersion, lock,. unlock, etc
func (c *QueryController) History(ctx *gin.Context) {
	var req mydomain.CellBase
	req, err := util.GetCellBaseFromParameter(ctx, false)
	if err != nil {
		msg := fmt.Sprintf("fail to parse parameter from url query, err:%v", err)
		ctx.JSON(http.StatusBadRequest, mydomain.NewErrorRespWithMsg(-1, msg))
		return
	}

	historyList, err := c.cellHistorySvc.Find(req.CellId, req.Branch)
	if err != nil {
		msg := fmt.Sprintf("fail to find history branch:%s, id:%d, err:%v", req.Branch, req.CellId, err)
		ctx.JSON(http.StatusInternalServerError, mydomain.NewErrorRespWithMsg(-1, msg))
	} else {
		msg := fmt.Sprintf("branch:%s, cellId:%d, size:%d", req.Branch, req.CellId, len(historyList))
		ctx.JSON(http.StatusOK, mydomain.NewErrorRespWithMsg(-1, msg))
	}

	return
}
