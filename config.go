/**
	Налаштовує роботу сервісу на підставі ini-файлу. 
	Приклад запуску: service -config ./service.ini
*/
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
	// Налаштування протоколювання
	// Рівень налагодження: debug, info, warn, error, dpanic, panic, fatal
	DebugLevel    = flag.String("DebugLevel", "debug", "")
	
	// Лог-файл, програма повинна мати доступ на запис до цього файлу
	LogPath       = flag.String("Log", "stdout", "")
	
	// Формат дати у лог-файлі
	LogTimeFormat = flag.String("LogTimeFormat", "02-01-2006 15:04:05", "")

	// Налаштування web service server
	// TCP-адреса сервера для прийому вхідних запитів
	Addr    = flag.String("Addr", "0.0.0.0:3000", "")
	
	// Таймаут
	Timeout = flag.Duration("Timeout", time.Second*15, "")
	
	// Джерело з якого дозволено запит до сервісу
	Origin  = flag.String("Origin", "0.0.0.0:4200", "")
	
	// Рівень налагодження gin: debug, realese, test
	GinMode = flag.String("GinMode", "debug", "")

	log *zap.SugaredLogger
)

// Кодер формату часу
func myTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(*LogTimeFormat))
}

func init() {
	iniflags.Parse()

	var atomicLevel zapcore.Level

	// Встановлюємо рівень налагодження
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

	// Конфігурація системи протоколювання 
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
