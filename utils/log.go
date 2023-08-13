package utils

import "github.com/sirupsen/logrus"

func InitLog() {
	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006/01/02 - 15:04:05",
	}

	// 设置 logrus 的默认格式化器
	logrus.SetFormatter(formatter)
}

func NewHandleLog(method string, handle string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"Method": method,
		"Handle": handle,
	})
}

func NewFuncLog(function string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"Function": function,
	})
}
