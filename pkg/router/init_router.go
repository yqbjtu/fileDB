package router

import (
	"fileDB/docs"
	mycontroller "fileDB/pkg/controller"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ConfigRouter(router *gin.Engine) {

	docs.SwaggerInfo.BasePath = "/api/v1"
	// use ginSwagger middleware to serve the API docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	cvsGroupEngine := router.Group("api/v1/cvs")
	cvsController := mycontroller.NewCvsController()
	cvsGroupEngine.POST("/add", cvsController.CreateNewVersion)
	cvsGroupEngine.POST("/lock", cvsController.Lock)
	cvsGroupEngine.POST("/unlock", cvsController.UnLock)

	queryGroupEngine := router.Group("api/v1/query")
	queryController := mycontroller.NewQueryController()
	queryGroupEngine.GET("/download", queryController.DownloadFile)
	queryGroupEngine.GET("/history", queryController.History)

	router.POST("/users/:userId", cvsController.UpdateOneUser)
	router.GET("/usersfind", cvsController.FindUsers)
	router.GET("/cellversion/status", cvsController.Status)

	router.GET("/users/:userId", cvsController.GetOneUser)

	miscGroupEngine := router.Group("api/v1/mics")
	miscController := mycontroller.NewMiscController()
	miscGroupEngine.GET("/freeMemory", miscController.FreeMemory)
}
