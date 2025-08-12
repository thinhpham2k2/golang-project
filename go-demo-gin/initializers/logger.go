package initializers

import (
	"io"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() {
	logPath := getEnv("LOG_FILE", "log/app.log")
	_ = os.MkdirAll(filepath.Dir(logPath), 0o755)

	rotator := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    50, // MB
		MaxBackups: 7,
		MaxAge:     30, // days
		Compress:   true,
	}

	// Nếu mở file/rotator lỗi thì vẫn có stdout
	log.SetOutput(io.MultiWriter(os.Stdout, rotator))

	// Dev: text; Prod: JSON (tuỳ env)
	if getEnv("LOG_FORMAT", "text") == "json" {
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	} else {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp:    true,
			DisableColors:    true, // vì có ghi ra file
			QuoteEmptyFields: true, // đảm bảo "source" có ngoặc kép
			ForceQuote:       true,
			TimestampFormat:  time.RFC3339Nano,
		})
	}

	// Level qua env: debug|info|warn|error
	switch getEnv("LOG_LEVEL", "info") {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	// (Tuỳ) in caller cho debug, tốn chút overhead
	// log.SetReportCaller(true)
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
