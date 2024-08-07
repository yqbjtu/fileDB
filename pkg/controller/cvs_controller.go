package controller

import (
	"fileDB/pkg/config"
	mydomain "fileDB/pkg/domain"
	"fileDB/pkg/service"
	"fileDB/pkg/store"
	"time"

	"fileDB/pkg/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/klog"
	"net/http"
	"strconv"
)

type CvsController struct {
	// service or some to access DB method
}

func NewCvsController() *CvsController {
	controller := CvsController{}
	return &controller
}

// CreateNewVersion 文件提交一个新版本，
func (c *CvsController) CreateNewVersion(ctx *gin.Context) {
	var req mydomain.AddVersionReq
	var err error
	cellIdStr := ctx.Query("cellId")
	versionStr := ctx.Query("version")
	branchStr := ctx.Query("branch")
	lockKeyStr := ctx.Query("lockKey")

	if cellIdStr == "" || branchStr == "" || versionStr == "" || lockKeyStr == "" {
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
				"errMsg": "version is not int type",
			})
			return
		}

		cellId, err := strconv.ParseInt(cellIdStr, 10, 32)
		if err != nil {
			klog.Errorf("failed to convert (%s)to int64, err:%v", cellIdStr, err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errMsg": "cellId is not int type",
			})
			return
		}

		req.CellId = cellId
		req.Branch = branchStr
		req.LockKey = lockKeyStr
	}

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
	baseOsmDataDir := config.GetConfig().OSMConfig.DataDir
	savePath := fmt.Sprintf("%s/%s/", baseOsmDataDir, req.Branch) + filename

	// 将上传的文件存储到服务器上指定的位置
	if err := ctx.SaveUploadedFile(header, savePath); err != nil {
		klog.Errorf("failed to write file %q, err:%v", filename, err)
		commentResult := mydomain.CommentResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to SaveUploadedFile, err:%v", err)}
		ctx.JSON(http.StatusOK, commentResult)
		return
	}

	var items []domain.CellStatus
	result := store.MyDB.Find(&items)
	fmt.Println("result:", result)
	cellStatusStore := store.NewCellStatusStore(store.MyDB)
	cellStatus, err := cellStatusStore.Find(req.CellId, branchStr)
	if err != nil {
		klog.Errorf("failed to find cell status, err:%v", err)
	}

	// if cell not exist, create a new cell status
	if cellStatus.CellId == 0 {
		cellStatus.CellId = req.CellId
		cellStatus.LatestVersion = req.Version
		cellStatus.LockKey = ""
		cellStatus.Branch = branchStr
		result = store.MyDB.Save(&cellStatus)
		if result.Error != nil {
			klog.Errorf("failed to save cell status, err:%v", result.Error)
			commentResult := mydomain.CommentResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to  save cell status, err:%v", result.Error)}
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
		errMsg := fmt.Sprintf("cellId:%s, current latest version is %d, expectedVersion should be %d, not %d", cellIdStr,
			cellStatus.LatestVersion, expectedVersion, req.Version)
		commentResult := mydomain.CommentResult{Code: -1, Data: nil, Msg: errMsg}
		ctx.JSON(http.StatusBadRequest, commentResult)
		return
	}

	// the cell should not be locked, or it is locked by req.LockKey
	if cellStatus.LockKey != "" && cellStatus.LockKey != req.LockKey {
		errMsg := fmt.Sprintf("cellId:%s is locked by %q, not %q", cellIdStr, cellStatus.LockKey, req.LockKey)
		commentResult := mydomain.CommentResult{Code: -1, Data: nil, Msg: errMsg}
		ctx.JSON(http.StatusBadRequest, commentResult)
		return
	}

	// update the cell status with latestVersion
	cellStatus.LatestVersion = req.Version
	cellStatus.LockKey = ""
	result = store.MyDB.Save(&cellStatus)
	if result.Error != nil {
		klog.Errorf("failed to save cell status, err:%v", result.Error)
		commentResult := mydomain.CommentResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to SaveUploadedFile, err:%v", result.Error)}
		ctx.JSON(http.StatusOK, commentResult)
		return
	}
	commentResult := mydomain.CommentResult{Code: 0, Data: req, Msg: "add new version ok"}

	cellHistoryStore := store.NewCellHistoryStore(store.MyDB)
	cellHistoryService := service.NewCellHistoryService(*cellHistoryStore)
	cellHistory := domain.CellHistory{
		CellId:      req.CellId,
		Branch:      req.Branch,
		Version:     req.Version,
		RequestType: "CheckinRequest",
		LockKey:     req.LockKey,
		Who:         "tester1",
	}
	cellHistoryService.Insert(cellHistory)
	ctx.JSON(http.StatusOK, commentResult)
}

func (c *CvsController) GetOneUser(context *gin.Context) {
	userId := context.Param("userId")
	klog.Infof("get one user by id %q", userId)

	context.JSON(http.StatusOK, gin.H{
		"searchId": userId,
	})
}

// cellId=1507888&branch=test
func (c *CvsController) Status(context *gin.Context) {
	cellIdStr := context.Query("cellId")
	branch := context.Query("branch")
	var cellId int64
	var err error
	if cellIdStr == "" {
		klog.Errorf("cellId can't be empty", cellIdStr)
		context.JSON(http.StatusBadRequest, gin.H{
			"errMsg": "cellId is empty",
		})
		return
	} else {
		cellId, err = strconv.ParseInt(cellIdStr, 10, 64)
		if err != nil {
			klog.Errorf("failed to convert (%s)to int64, err:%v", cellIdStr, err)
			context.JSON(http.StatusBadRequest, gin.H{
				"errMsg": "cellId is int type",
			})
			return
		}
	}

	if branch == "" {
		branch = "main"
		klog.Infof("get by cellId %v, use default branch %q", cellId, branch)
	} else {
		klog.Infof("get by cellId %v, branch %q", cellId, branch)
	}
	branches := [3]string{"main", "redo", "test"}

	response := map[string]interface{}{
		"version":  5,
		"cellId":   cellId,
		"branches": branches,
	}
	context.JSON(http.StatusOK, response)
}

func (c *CvsController) Lock(ctx *gin.Context) {
	// 从body中解析出cellId, plus1Ver, , branch
	var commentResult mydomain.CommentResult
	lockReq := mydomain.LockReq{}
	if err := ctx.ShouldBindJSON(&lockReq); err != nil {
		commentResult = mydomain.CommentResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to parse http body, err:%v", err)}
		ctx.JSON(http.StatusBadRequest, commentResult)
		return
	}

	if lockReq.CellId <= 0 {
		klog.Errorf("cellId can't be <= 0, req:%v", lockReq)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errMsg": fmt.Sprintf("cellId is %+v, it should be gte 0", lockReq),
		})
		return
	}

	if lockReq.Branch == "" {
		klog.Errorf("branch can't be <= 0, req:%v", lockReq)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errMsg": fmt.Sprintf("branch can't be empty"),
		})
		return
	}

	cellStatusStore := store.NewCellStatusStore(store.MyDB)
	cellStatus, err := cellStatusStore.Find(lockReq.CellId, lockReq.Branch)
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

	result := store.MyDB.Save(&cellStatus)
	if result.Error != nil {
		klog.Errorf("failed to save cell status, err:%v", result.Error)
		commentResult =
			mydomain.CommentResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to save cell status lock info, err:%v", result.Error)}
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

// UnLock 幂等操作， 也就是如果该cell没有被加锁，调用unlock会直接成功
// 给出英文注释

func (c *CvsController) UnLock(ctx *gin.Context) {
	var commentResult mydomain.CommentResult
	lockReq := mydomain.LockReq{}
	if err := ctx.ShouldBindJSON(&lockReq); err != nil {
		commentResult = mydomain.CommentResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to parse http body, err:%v", err)}
		ctx.JSON(http.StatusBadRequest, commentResult)
		return
	}

	if lockReq.CellId <= 0 {
		klog.Errorf("cellId can't be <= 0, req:%v", lockReq)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errMsg": fmt.Sprintf("cellId is %d, it should be > 0", lockReq.CellId),
		})
		return
	}

	if lockReq.Branch == "" {
		klog.Errorf("branch can't be <= 0, req:%v", lockReq)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errMsg": fmt.Sprintf("branch can't be empty"),
		})
		return
	}

	cellStatusStore := store.NewCellStatusStore(store.MyDB)
	cellStatus, err := cellStatusStore.Find(lockReq.CellId, lockReq.Branch)
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
	result := store.MyDB.Save(&cellStatus)
	if result.Error != nil {
		klog.Errorf("failed to save cell status, err:%v", result.Error)
		commentResult = mydomain.CommentResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to save cell status lock info, err:%v", result.Error)}
		ctx.JSON(http.StatusInternalServerError, commentResult)
		return

	}

	msgStr := fmt.Sprintf("cell %d is unlocked by %s now, unlock done",
		lockReq.CellId, cellStatus.LockKey)
	commentResult = mydomain.CommentResult{Code: -1, Data: nil, Msg: msgStr}
	ctx.JSON(http.StatusOK, commentResult)
	return

}

/*
// 匹配的url格式:  /usersfind?username=tom&email=test1@163.com
*/
func (c *CvsController) FindUsers(ctx *gin.Context) {
	userName := ctx.DefaultQuery("username", "张三")
	email := ctx.Query("email")
	// 执行实际搜索，这里只是示例
	ctx.String(http.StatusOK, "search user by %q %q", userName, email)
}

func (c *CvsController) UpdateOneUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	klog.Infof("update user by id %q", userId)
}

func (c *CvsController) DeleteOneUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	klog.Infof("delete user by id %q", userId)

}
