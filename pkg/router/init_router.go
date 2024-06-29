package router

import (
	mycontroller "fileDB/pkg/controller"
	"github.com/gin-gonic/gin"
)

func ConfigRouter(router *gin.Engine) {
	cvsGroupEngine := router.Group("api/v1/cvs")
	cvsController := mycontroller.NewCvsController()
	cvsGroupEngine.POST("/add", cvsController.CreateNewVersion)
	cvsGroupEngine.POST("/cellvcs/checkout", cvsController.Lock)

	router.POST("/users/:userId", cvsController.UpdateOneUser)
	//router.GET("/users", cvsController.GetAllUsers)
	router.GET("/usersfind", cvsController.FindUsers)
	router.GET("/cellversion/status", cvsController.Status)

	router.GET("/users/:userId", cvsController.GetOneUser)

	miscGroupEngine := router.Group("api/v1/mics")
	miscController := mycontroller.NewMiscController()
	miscGroupEngine.GET("/freeMemory", miscController.FreeMemory)
}
