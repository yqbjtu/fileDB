package util

import (
	mydomain "fileDB/pkg/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/klog"
	"strconv"
)

func GetCellBaseFromParameter(ctx *gin.Context, isWithVersion bool) (mydomain.CellBase, error) {
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
