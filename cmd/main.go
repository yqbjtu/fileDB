package main

import (
	"context"
	"fileDB/pkg/config"
	"fileDB/pkg/controller"
	"fileDB/pkg/log"
	myrouter "fileDB/pkg/router"
	"fileDB/pkg/service"
	"fileDB/pkg/store"
	"flag"
	"fmt"
	"go.uber.org/fx"
	"k8s.io/klog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ServiceLifetimeHooks -
func ServiceLifetimeHooks(lc fx.Lifecycle, srv *http.Server) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				fmt.Printf("starting web server listen and serve at %v ...", srv.Addr)
				go func() {
					if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						fmt.Printf("listen: %s\n", err)
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				fmt.Print("closing web server ...")
				ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
				defer cancel()
				return srv.Shutdown(ctx)
			},
		},
	)
}

/*
 the main function
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

	app := fx.New(
		fx.Provide(fx.Annotated{
			Name: "addr",
			Target: func() string {
				return fmt.Sprintf(":%d", config.GetConfig().Port)
			},
		}),
		fx.Provide(func() context.Context {
			return context.Background()
		}),
		config.Module,
		store.Module,
		service.Module,
		controller.Module,
		myrouter.Module,
		fx.Invoke(ServiceLifetimeHooks),
	)

	shutdowner := make(chan os.Signal, 1)
	signal.Notify(shutdowner, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGCHLD)

	go func() {
		if err := app.Start(context.Background()); err != nil {
			fmt.Printf("Failed to start the application: %v\n", err)
			os.Exit(1)
		}
	}()

	<-shutdowner

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := app.Stop(ctx); err != nil {
		fmt.Printf("Failed to gracefully stop the application: %v\n", err)
		os.Exit(2)
	}

	fmt.Println("Application stopped gracefully")
	//store.InitDB()
	//router := gin.Default()
	//myrouter.ConfigRouter(router)
	//
	//webServer := &http.Server{
	//	Addr:           ":8090",
	//	Handler:        router,
	//	ReadTimeout:    15 * time.Second,
	//	WriteTimeout:   15 * time.Second,
	//	MaxHeaderBytes: 1 << 20,
	//}
	//
	//webServer.ListenAndServe()

	//router.Run()
	// router.Run(":8090") 也能运行指定端口和ip上
}
