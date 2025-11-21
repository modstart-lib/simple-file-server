package module

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"simple-file-server/global"
	"simple-file-server/lib/files"
	"time"
)

var monitorCancel context.CancelFunc

type MonitorService struct {
}

func (m *MonitorService) Run() {
	log.Info("MonitorService Run")
	fileList := files.ListFiles(global.CONFIG.TempDir)
	for _, file := range fileList {
		if file.IsDir {
			continue
		}
		if file.Mtime < time.Now().Unix()-3600*24*30 {
			log.Info("CleanCacheFile:" + file.Name)
			files.DeleteFile(file.Path)
		}
	}
	// Clean MultiPart dirs older than 24 hours
	multiPartDir := global.CONFIG.TempDir + "/MultiPart"
	if files.FileExists(multiPartDir) {
		fileList = files.ListFiles(multiPartDir)
		for _, file := range fileList {
			if file.IsDir && file.Mtime < time.Now().Unix()-3600*24 {
				log.Info("CleanMultiPartDir:" + file.Name)
				files.DeleteDir(file.Path)
			}
		}
	}
}

func StartMonitor(removeBefore bool, interval int64) error {
	if removeBefore {
		monitorCancel()
		global.CRON.Remove(cron.EntryID(global.CronIDMonitor))
	}

	service := &MonitorService{}

	_, cancel := context.WithCancel(context.Background())
	monitorCancel = cancel
	monitorID, err := global.CRON.AddJob(fmt.Sprintf("@every %ds", interval), service)
	if err != nil {
		return err
	}
	global.CronIDMonitor = monitorID
	return nil
}
