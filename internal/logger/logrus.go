package logger

import "github.com/sirupsen/logrus"

func New() *logrus.Logger {
	log := logrus.New()

	// Включаем цвета (по умолчанию уже включены для TTY)
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})
	return log
}
