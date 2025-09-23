package logger

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

var globalLogger *slog.Logger
var logFileDir = "./runtime/log"

func init() {
	// 初始化默认日志
	globalLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: false,
		// 开销太大了，不记录上下文信息
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			// 格式化时间戳
			if attr.Key == slog.TimeKey {
				// 将时间格式化为字符串（使用日志记录的时间，即attr.Value）
				if t, ok := attr.Value.Any().(time.Time); ok {
					return slog.Attr{
						Key:   attr.Key,
						Value: slog.StringValue(t.Format(time.RFC3339)), // 或其他格式如"2006-01-02 15:04:05"
					}
				}
			}
			return attr
		},
	}))
	slog.SetDefault(globalLogger)
	// 用默认配置
}

func setLogFile() error {
	//确保存储目录存在
	if _, err := os.Stat(logFileDir); os.IsNotExist(err) {
		err = os.MkdirAll(logFileDir, 0777) //0777代表最宽松的权限
		if err != nil {
			return fmt.Errorf("create log dir '%s' error: %v", logFileDir, err)
		}
	}
	return nil
}

// 打开或创建文件
func writeToLogFile(filename string) (*os.File, error) {
	times := time.Now().Format(time.RFC3339)
	filePath := path.Join(logFileDir, fmt.Sprintf("%s_%s.log", filename, times))
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	//4=读权限，2=写权限，1=执行权限
	if err != nil {
		return nil, fmt.Errorf("fail to open file:%v", err)
	}
	return file, nil
}

// 将一定级别的日志写入文件
func logToFile(level slog.Level, logName string, fields ...map[string]any) error {
	file, err := writeToLogFile(logName)
	if err != nil {
		return err
	}
	defer file.Close()

	//创建并输出logger
	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{})
	fLogger := slog.New(handler)

	merge := map[string]any{}
	//字段合并
	for _, field := range fields {
		for key, value := range field {
			merge[key] = value
		}
	}

	//写入
	fLogger.LogAttrs(
		context.Background(),
		level,
		"log message",
		slog.Group("fields", slog.Any("", merge)),
	)
	return nil
}

// 这个函数是配置gin的日志，用结构化日志并且写入
// 原来的方法有错，注释掉了
func LoggerToFile() gin.HandlerFunc {
	// 按日期创建日志文件
	times := time.Now().Format("2006-01-02")
	filename := path.Join(logFileDir, fmt.Sprintf("%s.log", times))

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(fmt.Errorf("fail to open log: %v", err))
	}

	// 返回中间件
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 记录日志
		param := gin.LogFormatterParams{
			TimeStamp:    time.Now(),
			ClientIP:     c.ClientIP(),
			Method:       c.Request.Method,
			StatusCode:   c.Writer.Status(),
			Latency:      time.Since(start),
			Path:         path,
			ErrorMessage: c.Errors.ByType(gin.ErrorTypePrivate).String(),
		}

		if raw != "" {
			path = path + "?" + raw
		}

		message := fmt.Sprintf("%s - %s \"%s %s %d %s \"%s\" %s\"\n",
			param.TimeStamp.Format(time.RFC3339),
			param.ClientIP,
			param.Method,
			path,
			param.StatusCode,
			param.Latency,
			c.Request.UserAgent(),
			param.ErrorMessage,
		)

		_, _ = file.WriteString(message)
	}
}

// func LoggerToFile() gin.HandlerFunc {
// 	if err := setLogFile(); err != nil {
// 		panic(fmt.Errorf("fail to set log:%v", err))
// 	}
// 	times := time.Now().Format(time.RFC3339)
// 	filename := path.Join(logFileDir, fmt.Sprintf("%s.log", times))

// 	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
// 	if err != nil {
// 		panic(fmt.Errorf("fail to open log:%v", err))
// 	}

// 	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
// 		message := fmt.Sprintf("%s - %s \"%s %s %s %d %s \"%s\" %s\"\n",
// 			param.TimeStamp.Format(time.RFC3339), //日志时间戳
// 			param.ClientIP,                       //客户端ip
// 			param.Method,                         //http请求方法
// 			param.Path,                           //请求路径
// 			param.Request.Proto,
// 			//request是原始http请求对象，proto是协议版本
// 			param.StatusCode,          //状态码
// 			param.Latency,             //请求耗时
// 			param.Request.UserAgent(), //用户代理
// 			param.ErrorMessage,        //错误信息

// 		)
// 		_, _ = file.WriteString(message)
// 		return message
// 	})
// }

// 日志写入
func Write(msg string, filename string) {
	fields := map[string]any{
		"message": msg,
	}

	if err := logToFile(slog.LevelInfo, filename, fields); err != nil {
		fmt.Printf("fail to write log:%v", err)

	}
}

func Debug(fields map[string]any) {
	logToFile(slog.LevelDebug, "debug", fields)
}

func Info(fields map[string]any) {
	logToFile(slog.LevelInfo, "info", fields)
}

func Warn(fields map[string]any) {
	logToFile(slog.LevelWarn, "warn", fields)
}

func Error(fields map[string]any) {
	logToFile(slog.LevelError, "error", fields)
}

func Fatal(fields map[string]any) {
	logToFile(slog.LevelError, "fatal", fields)
	os.Exit(1)
}

func Panic(fields map[string]any) {
	logToFile(slog.LevelError, "panic", fields)
	panic(fields)
}

func Trace(fields map[string]any) {
	logToFile(slog.LevelDebug, "trace", fields)
}

// 捕获panic写入日志，相应错误
func Recover(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			times := time.Now().Format(time.RFC3339)
			filename := path.Join(logFileDir, fmt.Sprintf("error_%s.log", times))

			file, fileErr := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if fileErr != nil {
				fmt.Printf("fail to open error log:%v", fileErr)
				return
			}

			defer file.Close()

			stack := string(debug.Stack())
			logger := slog.New(slog.NewJSONHandler(file, nil))

			logger.Error(
				"recover from panic",
				slog.Any("error", err),
				slog.String("stacktrace", stack),
			)

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  fmt.Sprintf("%v", err),
			})
			ctx.Abort() //终止后续的中间件和处理函数
		}
	}()
	ctx.Next() //暂停当前中间件，执行后续中间件和处理函数
}
