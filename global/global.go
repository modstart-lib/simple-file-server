package global

import (
	"github.com/robfig/cron/v3"
	"simple-file-server/lib/defs"
)

var (
	CONFIG defs.Config
	CRON   *cron.Cron

	CronIDMonitor cron.EntryID
)
