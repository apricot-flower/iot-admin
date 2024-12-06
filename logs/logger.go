package logs

import (
	"os"
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.SugaredLogger

func init() {
	logPath := "./iot_admin_logs"
	fileName := "iot_admin"
	debugLogPath := path.Join(logPath, "debug_"+fileName+".log")
	infoLogPath := path.Join(logPath, "info_"+fileName+".log")
	warnLogPath := path.Join(logPath, "warn_"+fileName+".log")
	errorLogPath := path.Join(logPath, "error_"+fileName+".log")
	var coreArr []zapcore.Core
	//获取编码器
	encoderConfig := zap.NewProductionEncoderConfig()            //NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        //指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder //显示完整文件路径
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	//日志级别
	errPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //error级别
		return lev >= zap.ErrorLevel
	})
	debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //info和debug级别,debug级别是最低的
		return lev == zap.DebugLevel
	})
	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //info和debug级别,debug级别是最低的
		return lev == zap.InfoLevel
	})
	warnPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //info和debug级别,debug级别是最低的
		return lev == zap.WarnLevel
	})

	debugFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   debugLogPath, //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    10,           //文件大小限制,单位MB
		MaxBackups: 30,           //最大保留日志文件数量
		MaxAge:     30,           //日志文件保留天数
		Compress:   false,        //是否压缩处理
	})
	debugFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(debugFileWriteSyncer, zapcore.AddSync(os.Stdout)), debugPriority)

	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   infoLogPath, //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    10,          //文件大小限制,单位MB
		MaxBackups: 30,          //最大保留日志文件数量
		MaxAge:     30,          //日志文件保留天数
		Compress:   false,       //是否压缩处理
	})
	infoFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), infoPriority)

	warnFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   warnLogPath, //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    10,          //文件大小限制,单位MB
		MaxBackups: 30,          //最大保留日志文件数量
		MaxAge:     30,          //日志文件保留天数
		Compress:   false,       //是否压缩处理
	})
	warnFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(warnFileWriteSyncer, zapcore.AddSync(os.Stdout)), warnPriority)

	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   errorLogPath, //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    10,           //文件大小限制,单位MB
		MaxBackups: 30,           //最大保留日志文件数量
		MaxAge:     30,           //日志文件保留天数
		Compress:   false,        //是否压缩处理
	})
	errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), errPriority)

	coreArr = append(coreArr, infoFileCore)
	coreArr = append(coreArr, errorFileCore)
	coreArr = append(coreArr, debugFileCore)
	coreArr = append(coreArr, warnFileCore)
	log := zap.New(zapcore.NewTee(coreArr...), zap.AddCaller())
	_ = log.Sync()
	zap.AddCaller() //为显示文件名和行号，可省略
	Logger = log.Sugar()
}
