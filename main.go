package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/shisiying/go-blog/global"
	"github.com/shisiying/go-blog/internal/model"
	"github.com/shisiying/go-blog/internal/routers"
	"github.com/shisiying/go-blog/pkg/logger"
	"github.com/shisiying/go-blog/pkg/setting"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	port    string
	runMode string
	config  string

	isVersion bool

	buildTime    string
	buildVersion string
	gitCommitId  string
)

func init() {

	err := setupFlag()
	if err != nil {
		log.Fatalf("init.setupFlag err:%v", err)
	}

	err = setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err:%v", err)
	}

	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err:%v", err)
	}

	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDbEngine err: %v", err)
	}

}

func setupLogger() error {
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
		MaxSize:   600,
		MaxAge:    10,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)
	return nil
}

func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}

// @title blog
// @version 1.0
// @description Go-blog

func main() {
	if isVersion {
		fmt.Printf("build_time:%s\n", buildTime)
		fmt.Printf("build_version:%s\n", buildVersion)
		fmt.Printf("git_commit_id:%s\n", gitCommitId)
		return
	}
	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()

	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeOut,
		WriteTimeout:   global.ServerSetting.WriteTimeOut,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("s.ListenAndServe err: %v", err)
		}

		//等待中断信号
		quit := make(chan os.Signal)
		//接收 syscall.SIGINT 和sysycall.SIGTERM信号
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shuting down server...")

		//最大时间控制，通知该服务端它有5s的时间来处理原有的请求
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			log.Fatal("server forced to shutdown:", err)
		}
		log.Println("server exiting")
	}()

	s.ListenAndServe()
}

func setupSetting() error {
	s, err := setting.NewSetting(strings.Split(config, ",")...)

	if err != nil {
		return err
	}

	err = s.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}

	err = s.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}

	err = s.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}

	err = s.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		return err
	}

	err = s.ReadSection("Email", &global.EmailSetting)
	if err != nil {
		return err
	}

	global.ServerSetting.ReadTimeOut *= time.Second
	global.ServerSetting.WriteTimeOut *= time.Second
	global.JWTSetting.Expire *= time.Second

	if port != "" {
		global.ServerSetting.HttpPort = port
	}

	if runMode != "" {
		global.ServerSetting.RunMode = runMode
	}

	return nil

}

func setupFlag() error {
	flag.StringVar(&port, "port", "", "启动端口")
	flag.StringVar(&runMode, "mode", "", "启动模式")
	flag.StringVar(&config, "config", "", "指定要使用的配置文件路径")
	flag.BoolVar(&isVersion, "version", false, "编译信息")
	flag.Parse()

	return nil
}
