package log_wrapper

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/evendo-project/log-wrapper/writer_custom_log"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

const (
	DefaultLogPath          = "/logs"
	DefaultLifespanLogFile  = 5 * time.Hour * 24
	DefaultLifetimeLogFiles = 30 * time.Hour * 24
)

var logger *zap.Logger
var zapLogger *ZapLogger

func CreateDefaultLogger(logLevel string, nameService string, logPath string, printInConsole bool) {
	CreateLogger(logLevel, nameService, printInConsole, logPath, DefaultLifespanLogFile, DefaultLifetimeLogFiles)
}

func CreateLogger(logLevel string, nameService string, printInConsole bool, logPath string, lifespanLogFile time.Duration, lifetimeLogFiles time.Duration) {
	wc := writer_custom_log.New(logPath, nameService, printInConsole, lifespanLogFile, lifetimeLogFiles)

	encoderConfig := ecszap.NewDefaultEncoderConfig()
	level, _ := zapcore.ParseLevel(logLevel)

	core := ecszap.NewCore(encoderConfig, wc, level)
	logger = zap.New(core, zap.AddCaller())

	hostname, _ := os.Hostname()
	logger = logger.With(zap.String("service_name", nameService))
	logger = logger.With(zap.String("hostname", hostname))

	zapLogger = &ZapLogger{
		logger: logger,
	}
}

func GetLogger() zap.Logger {
	return *logger
}

func GetLoggerWithTraceId() (zap.Logger, string) {
	traceId := uuid.New().String()
	loggerWithTrace := logger.With(zap.String("trace_id", traceId))
	return *loggerWithTrace, traceId
}

func GetLoggerWithTraceIdGinPlug(c *gin.Context) (zap.Logger, string) {
	loggerWithTrace, traceId := GetLoggerWithTraceId()
	c.Header("X-Correlation-ID", traceId)
	return loggerWithTrace, traceId
}

func GetLoggerWithCustomTraceId(traceId string) zap.Logger {
	loggerWithTrace := logger.With(zap.String("trace_id", traceId))
	return *loggerWithTrace
}

func GetLoggerWithCustomTraceIdGinPlug(c *gin.Context, traceId string) zap.Logger {
	loggerWithTrace := GetLoggerWithCustomTraceId(traceId)
	c.Header("X-Correlation-ID", traceId)
	return loggerWithTrace
}

func GetRecoveryWithLoggerGin() gin.HandlerFunc {
	return ginzap.RecoveryWithZap(getMinimizedZapLogger(), true)
}

func getMinimizedZapLogger() ginzap.ZapLogger {
	return zapLogger
}
