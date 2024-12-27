package logger

import (
	"github.com/sirupsen/logrus"
)

// Создание логгера для настройки уровней логирования на основе переменных окружения
func InitLogger() *logrus.Logger {
	logger := logrus.New()

	//err := godotenv.Load()
	//if err != nil {
	//	logger.Warn("Error loading .env file, using default configuration")
	//}

	logLevel := "debug"
	switch logLevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	logFormat := "text"
	if logFormat == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	return logger
}
