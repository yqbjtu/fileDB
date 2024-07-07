package main

import (
	"fileDB/pkg/config"
	"fileDB/pkg/log"
	myrouter "fileDB/pkg/router"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/klog"
	"net/http"
	"os"
	"time"
)

/*
  GET http://127.0.0.1:8090/users
  GET http://127.0.0.1:8090/usersfind?username=tom&email=test1@163.com
  PUT http://127.0.0.1:8090/users  body {"userName":"tester1"}
  启动参数 --log_file=C:\F\ginportdemo.log --logtostderr=false --alsologtostderr=true
  --logtostderr=false表示输出到日志文件中，不再标准输出输出中展示，该参数默认值为true，
  --alsologtostderr[=false]: 同时输出日志到标准错误控制台和文件， 该参数为true后控制台和日志文件同时都有
  这是简单示例gin运行简单示例
*/

func main() {
	klog.InitFlags(nil)
	defer klog.Flush()

	flag.Parse()
	klog.Info("start gin webserver on specific port")

	cfgPath := "./conf/conf.yaml"
	absDir, err := os.Executable()
	if err != nil {
		fmt.Printf("start gin webserver on specific port from %s,"+
			" failed to get executable path: err:%s", cfgPath, err)
	} else {
		fmt.Printf("start gin webserver on specific port from %s, currDir:%s", cfgPath, absDir)
	}

	config.InitConfig(cfgPath)

	// init the log setting
	sugarLogger := log.InitLogger(&config.GetConfig().LogConfig)
	defer sugarLogger.Sync()

	//store.InitDB()
	router := gin.Default()
	myrouter.ConfigRouter(router)
	webServer := &http.Server{
		Addr:           ":8090",
		Handler:        router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	webServer.ListenAndServe()

	//router.Run()
	// router.Run(":8090") 也能运行指定端口和ip上
}
