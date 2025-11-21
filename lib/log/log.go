package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

func Init() {
	logrus.SetLevel(logrus.InfoLevel)
	loggerWriter := &lumberjack.Logger{
		Filename:   fmt.Sprintf("./log/%s.log", time.Now().Format("2006-01-02")),
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}
	multiWriter := io.MultiWriter(loggerWriter, os.Stdout)
	logrus.SetOutput(multiWriter)
}
