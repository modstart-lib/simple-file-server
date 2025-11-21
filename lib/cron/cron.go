package cron

import (
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"simple-file-server/global"
	"simple-file-server/lib/common"
	"simple-file-server/module"
	"time"
)

func Run() {
	nyc, _ := time.LoadLocation(common.LoadTimeZone())
	global.CRON = cron.New(
		cron.WithLocation(nyc),
		cron.WithChain(cron.Recover(cron.DefaultLogger)),
		cron.WithChain(cron.DelayIfStillRunning(cron.DefaultLogger)))
	if err := module.StartMonitor(false, 60*10); err != nil {
		log.Errorf("can not add monitor corn job: %s", err.Error())
	}
	global.CRON.Start()
}
