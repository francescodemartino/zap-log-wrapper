package log_wrapper

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
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

func GetRecoveryWithLoggerGin() gin.HandlerFunc {
	return ginzap.RecoveryWithZap(getMinimizedZapLogger(), true)
}

func getMinimizedZapLogger() ginzap.ZapLogger {
	return zapLogger
}
