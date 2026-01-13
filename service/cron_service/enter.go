package cron_service

import (
	"github.com/robfig/cron/v3"
	"time"
)

func Cron() {
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	crontab := cron.New(cron.WithSeconds(), cron.WithLocation(timezone))

	// 每天两点同步文章数据
	crontab.AddFunc("0 0 2 * * *", SyncArticle)

	crontab.Start()
}
