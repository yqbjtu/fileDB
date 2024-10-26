package controller

import (
	"errors"
	"fileDB/pkg/common"
	"fileDB/pkg/config"
	mydomain "fileDB/pkg/domain"
	"fileDB/pkg/log"
	"fileDB/pkg/service"
	"fileDB/pkg/store"
	"fileDB/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CvsController struct {
	// service or some to access DB method
	globalConfig    *config.GlobalConfig
	cellCvsSvc      *service.CellCvsService
	cellHistorySvc  *service.CellHistoryService
	cellStatusStore *store.CellStatusStore
}

func NewCvsController(globalConfig *config.GlobalConfig,
	cellCvsSvc *service.CellCvsService,
	cellHistorySvc *service.CellHistoryService,
	cellStatusStore *store.CellStatusStore) *CvsController {
	controller := CvsController{
		globalConfig:    globalConfig,
		cellCvsSvc:      cellCvsSvc,
		cellHistorySvc:  cellHistorySvc,
		cellStatusStore: cellStatusStore,
	}
	return &controller
}

// AddNewVersion 提交一个新版本，
// @Summary AddNewVersion 提交一个新版本，
// @Description submit a new version of cell
// @Tags query
// @Accept  json
// @Produce json
// @Success 200 {object} mydomain.CommonResult "ok"
// @Failure 400 {string} string "cellId,version and branch is required"
// @Router /api/v1/cvs/add [post]
func (c *CvsController) AddNewVersion(ctx *gin.Context) {
	var req mydomain.AddVersionReq
	var err error

	lockKeyStr := ctx.Query("lockKey")

	req.CellBase, err = util.GetCellBaseFromParameter(ctx, true)
	// lockKey can be empty when the cell is not locked by any key
	req.LockKey = lockKeyStr

	log.Infof("add new file version, req:%v", req)
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		msg := fmt.Sprintf("FormFile error: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, mydomain.NewErrorRespWithMsg(-1, msg))
		return
	}
	defer file.Close()

	// 你可以访问header来获取文件名称、文件大小和文件类型等信息
	filename := fmt.Sprintf("%d@@%s@@%d.osm", req.CellId, req.Branch, req.Version)
	// 定义文件保存路径
	baseOsmDataDir := c.globalConfig.OSMConfig.DataDir
	savePath := fmt.Sprintf("%s/%s/", baseOsmDataDir, req.Branch) + filename

	// 将上传的文件存储到服务器上指定的位置
	if err := ctx.SaveUploadedFile(header, savePath); err != nil {
		log.Errorf("failed to write file %q, err:%v", filename, err)

		msg := fmt.Sprintf("fail to SaveUploadedFile, err:%v", err)
		ctx.JSON(http.StatusOK, mydomain.NewErrorRespWithMsg(-1, msg))
		return
	}

	CommonResult, err := c.cellCvsSvc.AddNewVersion(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, CommonResult)
		return
	} else {
		ctx.JSON(http.StatusOK, CommonResult)
		return
	}
}

// @Summary lock the cell
// @Description lock the cell by cellId and branch
// @Tags query
// @Accept  json
// @Produce json
// @Success 200 {string} string	"ok"
// @Failure 400 {string} string "We need cellId and branch"
// @Router /api/v1/csv/lock [post]
func (c *CvsController) Lock(ctx *gin.Context) {
	// 从body中解析出cellId, plus1Ver, , branch
	var CommonResult mydomain.CommonResult
	lockReq, err := getLockUnLockReq(ctx)
	if err != nil {
		msg := fmt.Sprintf("fail to parse http body, err:%v", err)
		ctx.JSON(http.StatusBadRequest, mydomain.NewErrorRespWithMsg(-1, msg))
		return
	}

	log.Infof("lock  cell %d, branch %s", lockReq.CellId, lockReq.Branch)
	CommonResult, err = c.cellCvsSvc.Lock(&lockReq)
	if err != nil {
		if errors.Is(err, common.ErrDBOperationFailure) {
			ctx.JSON(http.StatusInternalServerError, CommonResult)
		} else {
			ctx.JSON(http.StatusBadRequest, CommonResult)
		}

	} else {
		ctx.JSON(http.StatusOK, CommonResult)
	}
	return
}

func getLockUnLockReq(ctx *gin.Context) (mydomain.LockReq, error) {
	lockReq := mydomain.LockReq{}
	if err := ctx.ShouldBindJSON(&lockReq); err != nil {
		return lockReq, fmt.Errorf("fail to parse http body, err:%v", err)
	}

	if lockReq.CellId <= 0 {
		return lockReq, fmt.Errorf("cellId is %d, it should be > 0", lockReq.CellId)
	}

	if lockReq.Branch == "" {
		log.Errorf("branch can't be <= 0, req:%v", lockReq)
		return lockReq, fmt.Errorf("branch can't be empty")
	}

	return lockReq, nil
}

func (c *CvsController) UnLock(ctx *gin.Context) {
	lockReq, err := getLockUnLockReq(ctx)
	if err != nil {
		msg := fmt.Sprintf("fail to parse http body, err:%v", err)
		ctx.JSON(http.StatusBadRequest, mydomain.NewErrorRespWithMsg(-1, msg))
		return
	}

	cellStatus, err := c.cellStatusStore.Find(lockReq.CellId, lockReq.Branch)
	if err != nil {
		log.Errorf("failed to find cell status, err:%v", err)
	}

	// if the cell is not locked, return ok
	if cellStatus.LockKey == "" {
		msg := fmt.Sprintf("cell %d has not been locked", lockReq.CellId)
		ctx.JSON(http.StatusOK, mydomain.NewSuccessRespWithMsg(nil, msg))
		return
	}

	if cellStatus.LockKey != lockReq.LockKey {
		// if the cell is locked by other lockKey, return fail
		msg := fmt.Sprintf("cell %d is locked by %s now, it can't be unlocked by %s",
			lockReq.CellId, cellStatus.LockKey, lockReq.LockKey)
		ctx.JSON(http.StatusConflict, mydomain.NewErrorRespWithMsg(-1, msg))
		return
	}

	// if the cell is locked by this lockKey, unlock and return ok
	cellStatus.LockKey = ""
	cellStatus.LockTimeFrom = nil
	cellStatus.LockTimeTo = nil
	_, err = c.cellStatusStore.Save(cellStatus)
	if err != nil {
		log.Errorf("failed to save cell status, err:%v", err)
		msg := fmt.Sprintf("fail to save cell status lock info, err:%v", err)
		ctx.JSON(http.StatusInternalServerError, mydomain.NewErrorRespWithMsg(-1, msg))
		return
	}

	msg := fmt.Sprintf("cell %d is unlocked by %s now, unlock done",
		lockReq.CellId, cellStatus.LockKey)
	ctx.JSON(http.StatusOK, mydomain.NewSuccessRespWithMsg(nil, msg))
	return
}
