package controller

import (
	"fileDB/pkg/config"
	mydomain "fileDB/pkg/domain"
	"fileDB/pkg/service"
	"fileDB/pkg/store"
	"fileDB/pkg/util"
	"time"

	"fileDB/pkg/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/klog"
	"net/http"
)

type CvsController struct {
	// service or some to access DB method
	globalConfig       *config.GlobalConfig
	cellHistoryService *service.CellHistoryService
	cellStatusStore    *store.CellStatusStore
}

func NewCvsController(globalConfig *config.GlobalConfig,
	cellHistoryService *service.CellHistoryService,
	cellStatusStore *store.CellStatusStore) *CvsController {
	controller := CvsController{
		globalConfig:       globalConfig,
		cellHistoryService: cellHistoryService,
		cellStatusStore:    cellStatusStore,
	}
	return &controller
}

// @Summary CreateNewVersion 提交一个新版本，
// @Description download cell file by cellId, version and branch
// @Tags query
// @Accept  json
// @Produce json
// @Success 200 {object} mydomain.CommentResult "ok"
// @Failure 400 {string} string "cellId,version and branch is required"
// @Router /api/v1/cvs/add[post]
func (c *CvsController) CreateNewVersion(ctx *gin.Context) {
	var req mydomain.AddVersionReq
	var err error

	lockKeyStr := ctx.Query("lockKey")

	req.CellBase, err = util.GetCellBaseFromParameter(ctx, true)
	// lockKey can be empty when the cell is not locked by any key
	req.LockKey = lockKeyStr

	klog.Infof("add new file version, req:%v", req)
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("FormFile error: %s", err.Error()),
		})
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
		klog.Errorf("failed to write file %q, err:%v", filename, err)
		commentResult := mydomain.CommentResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to SaveUploadedFile, err:%v", err)}
		ctx.JSON(http.StatusOK, commentResult)
		return
	}

	result, err := c.cellStatusStore.FindAll()
	fmt.Println("result:", result)

	cellStatus, err := c.cellStatusStore.Find(req.CellId, req.Branch)
	if err != nil {
		klog.Errorf("failed to find cell status, err:%v", err)
	}

	// if cell not exist, create a new cell status
	if cellStatus.CellId == 0 {
		cellStatus.CellId = req.CellId
		cellStatus.LatestVersion = req.Version
		cellStatus.LockKey = ""
		cellStatus.Branch = req.Branch
		_, err = c.cellStatusStore.Save(cellStatus)
		if err != nil {
			klog.Errorf("failed to save cell status, err:%v", err)
			commentResult := mydomain.CommentResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to  save cell status, err:%v", err)}
			ctx.JSON(http.StatusOK, commentResult)
			return
		} else {
			commentResult := mydomain.CommentResult{Code: 0, Data: req, Msg: "add the first version ok"}
			ctx.JSON(http.StatusOK, commentResult)
			return
		}
	}

	// the req.Version should be the latest version + 1
	expectedVersion := cellStatus.LatestVersion + 1
	if req.Version != expectedVersion {
		errMsg := fmt.Sprintf("cellId:%d, current latest version is %d, expectedVersion should be %d, not %d", req.CellId,
			cellStatus.LatestVersion, expectedVersion, req.Version)
		commentResult := mydomain.CommentResult{Code: -1, Data: nil, Msg: errMsg}
		ctx.JSON(http.StatusBadRequest, commentResult)
		return
	}

	// the cell should not be locked, or it is locked by req.LockKey
	if cellStatus.LockKey != "" && cellStatus.LockKey != req.LockKey {
		errMsg := fmt.Sprintf("cellId:%d is locked by %q, not %q", req.CellId, cellStatus.LockKey, req.LockKey)
		commentResult := mydomain.CommentResult{Code: -1, Data: nil, Msg: errMsg}
		ctx.JSON(http.StatusBadRequest, commentResult)
		return
	}

	// update the cell status with latestVersion
	cellStatus.LatestVersion = req.Version
	cellStatus.LockKey = ""
	_, err = c.cellStatusStore.Save(cellStatus)
	if err != nil {
		klog.Errorf("failed to save cell status, err:%v", err)
		commentResult := mydomain.CommentResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to SaveUploadedFile, err:%v", err)}
		ctx.JSON(http.StatusOK, commentResult)
		return
	}
	commentResult := mydomain.CommentResult{Code: 0, Data: req, Msg: "add new version ok"}

	cellHistory := domain.CellHistory{
		CellId:      req.CellId,
		Branch:      req.Branch,
		Version:     req.Version,
		RequestType: "CheckinRequest",
		LockKey:     req.LockKey,
		Who:         "tester1",
	}

	c.cellHistoryService.Insert(cellHistory)
	ctx.JSON(http.StatusOK, commentResult)
}

func (c *CvsController) GetOneUser(context *gin.Context) {
	userId := context.Param("userId")
	klog.Infof("get one user by id %q", userId)

	context.JSON(http.StatusOK, gin.H{
		"searchId": userId,
	})
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
	var commentResult mydomain.CommentResult
	lockReq, err := getLockUnLockReq(ctx)
	if err != nil {
		commentResult = mydomain.CommentResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to parse http body, err:%v", err)}
		ctx.JSON(http.StatusBadRequest, commentResult)
		return
	}

	cellStatus, err := c.cellStatusStore.Find(lockReq.CellId, lockReq.Branch)
	if err != nil {
		klog.Errorf("failed to find cell status, err:%v", err)
	}

	if cellStatus.LockKey != "" && cellStatus.LockKey != lockReq.LockKey {
		errMsg := fmt.Sprintf("cell %d has already been locked by %s now, so it can't be locked by %s again",
			lockReq.CellId, cellStatus.LockKey, lockReq.LockKey)
		commentResult = mydomain.CommentResult{Code: -1, Data: nil, Msg: errMsg}
		ctx.JSON(http.StatusBadRequest, commentResult)
		return
	}
	if lockReq.LockDuration.GetSeconds() <= 10 {
		errMsg := fmt.Sprintf("cell lock duration should be gt 10s, but it is %v", lockReq.LockDuration)
		commentResult = mydomain.CommentResult{Code: -1, Data: nil, Msg: errMsg}
		ctx.JSON(http.StatusBadRequest, commentResult)
		return
	}

	cellStatus.LockKey = lockReq.LockKey
	// cellStatus.LockTimeFrom等于当前时间
	fromTime := time.Now()
	cellStatus.LockTimeFrom = &fromTime

	// cellStatus.LockTimeTo等于当前时间加上一个小时
	goDuration := time.Duration(lockReq.LockDuration.GetSeconds())*time.Second + time.Duration(lockReq.LockDuration.GetNanos())*time.Nanosecond
	toTime := time.Now().Add(goDuration)
	cellStatus.LockTimeTo = &toTime
	if cellStatus.CellId == 0 {
		cellStatus.CellId = lockReq.CellId
		cellStatus.Branch = lockReq.Branch
		// 没有添加过，版本就为0
		cellStatus.LatestVersion = 0
	}

	_, err = c.cellStatusStore.Save(cellStatus)
	if err != nil {
		klog.Errorf("failed to save cell status, err:%v", err)
		commentResult =
			mydomain.CommentResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to save cell status lock info, err:%v", err)}
		ctx.JSON(http.StatusInternalServerError, commentResult)
		return
	}
	customTimeFormat := "2006-01-02 15:04:05"

	// add lock record in db
	response := map[string]interface{}{
		"id":           lockReq.CellId,
		"latestVer":    cellStatus.LatestVersion,
		"branch":       lockReq.Branch,
		"lockKey":      lockReq.LockKey,
		"lockTimeFrom": cellStatus.LockTimeFrom.Format(customTimeFormat),
		"lockTimeTo":   cellStatus.LockTimeTo.Format(customTimeFormat),
	}

	commentResult = mydomain.CommentResult{Code: 0, Data: response, Msg: "success"}
	ctx.JSON(http.StatusOK, response)
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
		klog.Errorf("branch can't be <= 0, req:%v", lockReq)
		return lockReq, fmt.Errorf("branch can't be empty")
	}

	return lockReq, nil
}

func (c *CvsController) UnLock(ctx *gin.Context) {
	var commentResult mydomain.CommentResult
	lockReq, err := getLockUnLockReq(ctx)
	if err != nil {
		commentResult = mydomain.CommentResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to parse http body, err:%v", err)}
		ctx.JSON(http.StatusBadRequest, commentResult)
		return
	}

	cellStatus, err := c.cellStatusStore.Find(lockReq.CellId, lockReq.Branch)
	if err != nil {
		klog.Errorf("failed to find cell status, err:%v", err)
	}

	// if the cell is not locked, return ok
	if cellStatus.LockKey == "" {
		msgStr := fmt.Sprintf("cell %d has not been locked", lockReq.CellId)
		commentResult = mydomain.CommentResult{Code: 0, Data: nil, Msg: msgStr}
		ctx.JSON(http.StatusOK, commentResult)
		return
	}

	if cellStatus.LockKey != lockReq.LockKey {
		// if the cell is locked by other lockKey, return fail
		msgStr := fmt.Sprintf("cell %d is locked by %s now, it can't be unlocked by %s",
			lockReq.CellId, cellStatus.LockKey, lockReq.LockKey)
		commentResult = mydomain.CommentResult{Code: -1, Data: nil, Msg: msgStr}
		ctx.JSON(http.StatusConflict, commentResult)
		return
	}

	// if the cell is locked by this lockKey, unlock and return ok
	cellStatus.LockKey = ""
	cellStatus.LockTimeFrom = nil
	cellStatus.LockTimeTo = nil
	_, err = c.cellStatusStore.Save(cellStatus)
	if err != nil {
		klog.Errorf("failed to save cell status, err:%v", err)
		commentResult = mydomain.CommentResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to save cell status lock info, err:%v", err)}
		ctx.JSON(http.StatusInternalServerError, commentResult)
		return

	}

	msgStr := fmt.Sprintf("cell %d is unlocked by %s now, unlock done",
		lockReq.CellId, cellStatus.LockKey)
	commentResult = mydomain.CommentResult{Code: -1, Data: nil, Msg: msgStr}
	ctx.JSON(http.StatusOK, commentResult)
	return
}
