package process

import "github.com/fumiama/cron"

// CronTab 全局定时器
var CronTab = cron.New()

func init() {
	CronTab.Start()
}
