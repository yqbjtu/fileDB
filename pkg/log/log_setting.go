package log

import (
	"fileDB/pkg/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger

func InitLogger(logConfig *config.LogConfig) *zap.SugaredLogger {
	writeSyncer := getLogWriter(logConfig)
	encoder := getEncoder()

	//   配置文件中读取level， 通过switch转换为zapcore中对应level
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	//logger := zap.New(core)
	//添加将调用函数信息记录到日志中的功能
	logger := zap.New(core, zap.AddCaller())

	//  根据配置文件配置是否显示caller信息
	logger.WithOptions(zap.AddCallerSkip(1))
	//logger.WithOptions(zap.WithCaller(true))

	Log = logger.Sugar()
	return Log
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()

	// 修改时间编码器
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 在日志文件中使用大写字母记录日志级别
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 使用console格式就失去了json 字段名称部分信息
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(logConfig *config.LogConfig) zapcore.WriteSyncer {
	//file, _ := os.Create(logConfig.LogPath)

	loggerFile := lumberjack.Logger{
		Filename: logConfig.LogPath,
		// 日志文件每1MB会切割并且在当前目录下最多保存5个备份
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
		//参数含义
		//Filename: 日志文件的位置
		//MaxSize：在进行切割之前，日志文件的最大大小（以MB为单位）
		//MaxBackups：保留旧文件的最大个数
		//MaxAges：保留旧文件的最大天数
		//Compress：是否压缩/归档旧文件
	}

	return zapcore.AddSync(&loggerFile)
}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(args ...interface{}) {
	Log.Debug(args)
}

// Info uses fmt.Sprint to construct and log a message.
func Info(args ...interface{}) {
	Log.Info(args)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) {
	Log.Warn(args)
}

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) {
	Log.Error(args)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(args ...interface{}) {
	Log.Panic(args)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	Log.Fatal(args)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	Log.Debugf(template, args)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	Log.Infof(template, args)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	Log.Warnf(template, args)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	Log.Errorf(template, args)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics. (See DPanicLevel for details.)
func DPanicf(template string, args ...interface{}) {
	Log.DPanicf(template, args)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(template string, args ...interface{}) {
	Log.Panicf(template, args)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(template string, args ...interface{}) {
	Log.Fatalf(template, args)
}
