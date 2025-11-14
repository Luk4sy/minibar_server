package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/service/log_service"
)

func main() {
	flags.Parse()                   // 1. 解析命令行参数：flags
	global.Config = core.ReadConf() // 2. 加载配置文件：conf + core/init_conf.go
	core.InitLogrus()               // 3. 初始化日志：core/init_logrus.go + logs
	global.DB = core.InitDB()       // 4. 初始化数据库：core/init_db.go

	log := log_service.NewRuntimeLog("同步文章数据", log_service.RuntimeDateHour)
	log.SetItem("文章1", 11)
	log.Save()
	log.SetItem("文章2", 12)
	log.Save()
}
