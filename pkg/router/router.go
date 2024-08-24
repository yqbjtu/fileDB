package router

import (
	mycontroller "fileDB/pkg/controller"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
)

// Router router
type Router struct {
	cvsController   *mycontroller.CvsController
	miscController  *mycontroller.MiscController
	queryController *mycontroller.QueryController
}

// NewRouter Generator
func NewRouter(
	cvsController *mycontroller.CvsController,
	miscController *mycontroller.MiscController,
	queryController *mycontroller.QueryController) *Router {
	return &Router{
		cvsController:   cvsController,
		miscController:  miscController,
		queryController: queryController,
	}
}

// ServerOption fx需要
type ServerOption struct {
	fx.In
	Addr       string            `name:"addr"`
	Middleware []gin.HandlerFunc `group:"middleware"`
}

// NewHTTPServer fx需要
func NewHTTPServer(router *Router, option ServerOption) *http.Server {
	return &http.Server{
		Addr:    option.Addr,
		Handler: router.Server(option.Middleware...),
	}
}

// Server main server
func (r *Router) Server(middlewares ...gin.HandlerFunc) *gin.Engine {
	gin.DisableConsoleColor()

	e := gin.New()
	// Setup middlewares
	e.Use(middlewares...)
	// Api router
	e.GET("", r.miscController.BuildInfo)
	{
		prefix := "/api/v1"
		baseEngine := e.Group(prefix)

		{
			cvsGroupEngine := e.Group(baseEngine.BasePath() + "/cvs")
			cvsGroupEngine.POST("/add", r.cvsController.AddNewVersion)
			cvsGroupEngine.POST("/lock", r.cvsController.Lock)
			cvsGroupEngine.POST("/unlock", r.cvsController.UnLock)
		}
		{
			miscGroupEngine := e.Group(baseEngine.BasePath() + "/mics")
			miscGroupEngine.GET("/freeMemory", r.miscController.FreeMemory)
			miscGroupEngine.GET("/buildInfo", r.miscController.BuildInfo)
		}
		{
			queryGroupEngine := e.Group(baseEngine.BasePath() + "/query")
			queryGroupEngine.GET("/download", r.queryController.DownloadFile)
			queryGroupEngine.GET("/history", r.queryController.History)
			queryGroupEngine.GET("/status", r.queryController.CellStatus)
		}
	}

	return e
}
