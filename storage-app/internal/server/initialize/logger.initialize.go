package initialize

import (
	"io"
	"os"

	"github.com/baothaihcmut/Bibox/storage-app/internal/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitializeLogger(cfg *config.LoggerConfig) *logrus.Logger {
	l := logrus.New()
	//parse level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	l.SetLevel(level)
	// set formater
	if cfg.Format == "json" {
		l.SetFormatter(&logrus.JSONFormatter{})
	} else {
		l.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}
	//set ouput
	var writers []io.Writer

	if cfg.Output == "stdout" || cfg.Output == "both" {
		writers = append(writers, os.Stdout)
	}

	if cfg.Output == "file" || cfg.Output == "both" {
		fileLogger := &lumberjack.Logger{
			Filename:   cfg.FileConfig.Path,
			MaxSize:    cfg.FileConfig.MaxSize,
			MaxBackups: cfg.FileConfig.MaxBackups,
			MaxAge:     cfg.FileConfig.MaxAge,
			Compress:   cfg.FileConfig.Compress,
		}
		writers = append(writers, fileLogger)
	}

	// Combine writers
	if len(writers) > 0 {
		l.SetOutput(io.MultiWriter(writers...))
	}
	return l

}
