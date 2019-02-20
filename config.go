package server_template

import (
	"flag"
	"github.com/vharitonsky/iniflags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"
)

var (
	// set logs
	DebugLevel    = flag.String("DebugLevel", "debug", "")
	LogPath       = flag.String("Log", "stdout", "")
	LogTimeFormat = flag.String("LogTimeFormat", "02-01-2006 15:04:05", "")

	// set web service server
	Addr    = flag.String("Addr", "0.0.0.0:3000", "")
	Timeout = flag.Duration("Timeout", time.Second*15, "")
	Origin  = flag.String("Origin", "0.0.0.0:4200", "")
	GinMode = flag.String("GinMode", "debug", "")

	log *zap.SugaredLogger
)

func myTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(*LogTimeFormat))
}

func init() {
	iniflags.Parse()

	// set debug level
	var atomicLevel zapcore.Level

	switch *DebugLevel {
	case "debug":
		atomicLevel = zapcore.DebugLevel
	case "info":
		atomicLevel = zapcore.InfoLevel
	case "warn":
		atomicLevel = zapcore.WarnLevel
	case "error":
		atomicLevel = zapcore.ErrorLevel
	case "dpanic":
		atomicLevel = zapcore.DPanicLevel
	case "panic":
		atomicLevel = zapcore.PanicLevel
	case "fatal":
		atomicLevel = zapcore.FatalLevel
	default:
		atomicLevel = zapcore.InfoLevel
	}

	// Set loggin systems
	cfg := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(atomicLevel),
		OutputPaths:      strings.Split(*LogPath, ","),
		ErrorOutputPaths: strings.Split(*LogPath, ","),
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: myTimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, _ := cfg.Build()
	defer logger.Sync() // flushes buffer, if any
	log = logger.Sugar()

	log.Warnf("Debug level: %v", atomicLevel)
	log.Debugf("Log: %v", *LogPath)
	log.Debugf("LogTimeFormat: %v", *LogTimeFormat)

	log.Debugf("Addr: %v", *Addr)
	log.Debugf("Timeout: %v", *Timeout)
	log.Debugf("Origin: %v", *Origin)
	log.Debugf("GinMode: %v", *GinMode)
}
