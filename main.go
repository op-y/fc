package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/op-y/fc/application/controller"
    "github.com/op-y/fc/config"

    "github.com/gin-gonic/gin"
     rotatelogs "github.com/lestrrat-go/file-rotatelogs"
    "github.com/rifflock/lfshook"
    "github.com/sirupsen/logrus"
)

func main() {
    // load config
    var cf string
    flag.StringVar(&cf, "c", "conf/fc.yaml", "fc config file path")
    flag.Parse()
    config.Cfg = config.LoadConfig(cf)

    // init logger
    rl, _ := rotatelogs.New(
        config.Cfg.AppLog + "%Y%m%d%H",
        rotatelogs.WithLinkName(config.Cfg.AppLog),
        rotatelogs.WithRotationTime(time.Hour),
        rotatelogs.WithMaxAge(time.Duration(3) * time.Hour),
    )
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    log.SetOutput(rl)
    log.Printf("===fc startup===")
    log.Printf("config file path is :%s", cf)
    config.Cfg.Print()

    // gin
    log.Printf("gin will start with port:%s", config.Cfg.Port)
    gin.SetMode(gin.ReleaseMode)
    gin.DisableConsoleColor()
    routes := gin.Default()
    routes.Use(logerMiddleware())
    go controller.StartGin(config.Cfg.Port, routes)
    //f, _ := os.Create("gin.log")
    //gin.DefaultWriter = io.MultiWriter(f)
    //gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

    // waiting
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigs
        fmt.Println()
        os.Exit(0)
    }()
    select {}
}

func logerMiddleware() gin.HandlerFunc {
	src, err := os.OpenFile(config.Cfg.GinLog, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Printf("failed to open gin log file: %s", err.Error())
	}
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.Out = src

	logWriter, err := rotatelogs.New(
		config.Cfg.GinLog + ".%Y%m%d%H.log",
		rotatelogs.WithLinkName(config.Cfg.GinLog),
		rotatelogs.WithMaxAge(time.Duration(3) * time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	logger.AddHook(lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}))

	return func(c *gin.Context) {
		//开始时间
		startTime := time.Now()
		//处理请求
		c.Next()
		//结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		//请求方式
		reqMethod := c.Request.Method
		//请求路由
		reqUrl := c.Request.RequestURI
		//状态码
		statusCode := c.Writer.Status()
		//请求IP
		clientIP := c.ClientIP()
		// 日志格式
		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUrl,
		}).Info()
	}
}
