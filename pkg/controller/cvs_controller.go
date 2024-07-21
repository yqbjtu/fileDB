package controller

import (
	"fileDB/pkg/config"
	mydomain "fileDB/pkg/domain"
	"fileDB/pkg/store"

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
				"errMsg": "version is not int",
			})
			return
		}

		req.CellId = cellIdStr
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
	filename := fmt.Sprintf("%s@@%s@@%d.osm", req.CellId, req.Branch, req.Version)
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
	cellStatus, err := cellStatusStore.Find(cellIdStr, branchStr)
	if err != nil {
		klog.Errorf("failed to find cell status, err:%v", err)
	}

	// if cell not exist, create a new cell status
	if cellStatus.CellId == 0 {
		cellId, _ := strconv.ParseInt(cellIdStr, 10, 32)
		cellStatus.CellId = cellId
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
				"errMsg": "cellId is int",
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
	if err := ctx.ShouldBind(&lockReq); err != nil {
		commentResult = mydomain.CommentResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to parse http body, err:%v", err)}
		ctx.JSON(http.StatusBadRequest, commentResult)
		return
	}

	if lockReq.CellId == "" {
		klog.Errorf("cellId can't be empty, req:%v", lockReq)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errMsg": "cellId is empty",
		})
		return
	}

	cellStatusStore := store.NewCellStatusStore(store.MyDB)
	cellStatus, err := cellStatusStore.Find(lockReq.CellId, lockReq.Branch)
	if err != nil {
		klog.Errorf("failed to find cell status, err:%v", err)
	}

	if cellStatus.LockKey != "" && cellStatus.LockKey != lockReq.LockKey {
		errMsg := fmt.Sprintf("cell %s has already locked by %s now, so it can't be locked by %s again",
			lockReq.CellId, cellStatus.LockKey, lockReq.LockKey)
		commentResult = mydomain.CommentResult{Code: -1, Data: nil, Msg: errMsg}
		ctx.JSON(http.StatusBadRequest, commentResult)
		return
	}

	// add lock record in db
	response := map[string]interface{}{
		"id": lockReq.CellId,
		//
		//"latestVer":       lockReq.Version,
		"ns":        lockReq.Branch,
		"lockStart": "",
		"lockEnd":   "",
		"lockKey":   lockReq.LockKey,
	}

	commentResult = mydomain.CommentResult{Code: 0, Data: response, Msg: "success"}
	ctx.JSON(http.StatusOK, response)
}

// UnLock 幂等操作， 也就是如果该cell没有被加锁，调用unlock会直接成功
// 给出英文注释

func (c *CvsController) UnLock(ctx *gin.Context) {
	var commentResult mydomain.CommentResult
	lockReq := mydomain.LockReq{}
	if err := ctx.ShouldBind(&lockReq); err != nil {
		commentResult = mydomain.CommentResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to parse http body, err:%v", err)}
		ctx.JSON(http.StatusBadRequest, commentResult)
		return
	}

	if lockReq.CellId == "" {
		klog.Errorf("cellId can't be empty", lockReq.CellId)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errMsg": "cellId is empty",
		})
		return
	}

	// if the cell is not locked, return ok
	// if the cell is locked by this lockKey, unlock and return ok
	// if the cell is locked by other lockKey, return fail

	ctx.JSON(http.StatusOK, "response")
}

/*
// 匹配的url格式:  /usersfind?username=tom&email=test1@163.com
*/
func (c *CvsController) FindUsers(context *gin.Context) {
	userName := context.DefaultQuery("username", "张三")
	email := context.Query("email")
	// 执行实际搜索，这里只是示例
	context.String(http.StatusOK, "search user by %q %q", userName, email)
}

func (c *CvsController) UpdateOneUser(context *gin.Context) {
	userId := context.Param("userId")
	klog.Infof("update user by id %q", userId)
}

func (c *CvsController) DeleteOneUser(context *gin.Context) {
	userId := context.Param("userId")
	klog.Infof("delete user by id %q", userId)

}
