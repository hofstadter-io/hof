package runtime

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (RT *Runtime) setupLogger(lvl string) {
	RT.setLogLevel(lvl)
	fmt.Println("SET LOGGER", lvl, RT.params.LogLevel)

	// filter below a level
	filter := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		fmt.Println("UGH: ", lvl, RT.params.LogLevel)
		return lvl >= RT.params.LogLevel
	})

	// onsole output
	console := zapcore.Lock(RT.stdout.(zapcore.WriteSyncer))

	// setup our config and console encoder
	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("15:04:05.000"))
	}

	consoleEncoder := zapcore.NewConsoleEncoder(config)

	core := zapcore.NewCore(consoleEncoder, console, filter)
	logger := zap.New(core)
	RT.logger = logger.Sugar()
}

func (RT *Runtime) setLogLevel(lvl string) {
	level := zapcore.ErrorLevel

	switch strings.ToLower(lvl) {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		panic("invalid log level: "+ lvl )
	}

	RT.params.LogLevel = level
}
