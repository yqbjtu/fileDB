package controller

import (
	mydomain "fileDB/pkg/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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
	namespaceStr := ctx.Query("namespace")
	lockKeyStr := ctx.Query("lockKey")

	if cellIdStr == "" || namespaceStr == "" || versionStr == "" || lockKeyStr == "" {
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
	filename := fmt.Sprintf("%s@@%s@@%d.osm", req.CellId, req.Namespace, req.Version)
	// 定义文件保存路径
	savePath := fmt.Sprintf("/tmp/osmdb/data/%s/", req.Namespace) + filename

	// 将上传的文件存储到服务器上指定的位置
	if err := ctx.SaveUploadedFile(header, savePath); err != nil {
		klog.Errorf("failed to write file %q, err:%v", filename, err)
		commentResult := mydomain.CommentResult{Code: -1, Data: nil, Msg: fmt.Sprintf("fail to SaveUploadedFile, err:%v", err)}
		ctx.JSON(http.StatusOK, commentResult)
		return
	}

	commentResult := mydomain.CommentResult{Code: 0, Data: req, Msg: "add new version ok"}
	ctx.JSON(http.StatusOK, commentResult)
}

//func (c *CvsController) GetAllUsers(context *gin.Context) {
//	klog.Infof("get all user")
//	//H is a shortcut for map[string]interface{}
//
//	var users []mydomain.User
//	var i int64
//	for i = 0; i < 3; i++ {
//		userName := fmt.Sprintf("tom%d", i)
//		user := mydomain.User{UserId: i}
//		user.UserName = userName
//		users = append(users, user)
//	}
//
//	context.JSON(http.StatusOK, gin.H{
//		"result": users,
//		"count":  len(users),
//	})
//}

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

// cvs/lock?cellId=1507888&namespace=test&plus1Ver=2&cellFilePath=%2F
func (c *CvsController) Lock(context *gin.Context) {
	cellIdStr := context.Query("cellId")
	cellFilePath := context.Query("cellFilePath")
	plus1VerStr := context.Query("plus1Ver")
	branch := context.DefaultQuery("branch", "main")

	var cellId int64
	var plus1Ver int64
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
				"errMsg": "cellId is not int",
			})
			return
		}
	}

	if plus1VerStr == "" {
		klog.Errorf("plus1Ver can't be empty", plus1VerStr)
		context.JSON(http.StatusBadRequest, gin.H{
			"errMsg": "plus1Ver is empty",
		})
		return
	} else {
		plus1Ver, err = strconv.ParseInt(plus1VerStr, 10, 64)
		if err != nil {
			klog.Errorf("failed to convert (%s)to int64, err:%v", plus1VerStr, err)
			context.JSON(http.StatusBadRequest, gin.H{
				"errMsg": "plus1Ver is not int",
			})
			return
		}
	}

	if cellFilePath == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"errMsg": "cellFilePath is empty",
		})
		return
	}

	body, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		klog.Infof("failed to read body, err:%v", err)
	} else {
		klog.Infof("body:%v", body)
	}

	request := map[string]string{"who": "user1", "jobId": "algo.work", "bizVersion": "1.1.0"}
	response := map[string]interface{}{
		"ver":       plus1Ver,
		"id":        cellId,
		"branch":    branch,
		"timestamp": "",
		"request":   request,
	}
	context.JSON(http.StatusOK, response)
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
