package log_service

import (
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
)

type Runtimelog struct {
	level           enum.LogLevelType
	title           string
	itemList        []string
	serviceName     string
	runtimeDateType RuntimeDateType
}

type RuntimeDateType int8

const (
	RuntimeDateHour  RuntimeDateType = 1
	RuntimeDateDay   RuntimeDateType = 2
	RuntimeDateWeek  RuntimeDateType = 3
	RuntimeDateMonth RuntimeDateType = 4
)

func (r *Runtimelog) Save() {

	// 判断创建还是更新
	var query gorm.DB
	switch r.runtimeDateType {
	case RuntimeDateHour:
	case RuntimeDateDay:
		query.Where("date (created_at) = date(now())")
	case RuntimeDateWeek:
	case RuntimeDateMonth:

	}

	var log models.LogModel
	global.DB.Where(query).Find(&log, "service_name = ? and log_type = ? and ", r.serviceName, enum.RuntimeLogType)
	content := strings.Join(r.itemList, "\n")

	if log.ID != 0 {
		// 更新
		return
	}
	err := global.DB.Create(&models.LogModel{
		LogType: enum.RuntimeLogType,
		Title:   r.title,
		Content: content,
		Level:   r.level,
	}).Error
	if err != nil {
		logrus.Errorf("创建运行日志错误 %s", err)
		return
	}

}

func NewRuntimeLog(serviceName string, dateType RuntimeDateType) *Runtimelog {
	return &Runtimelog{
		serviceName:     serviceName,
		runtimeDateType: dateType,
	}
}
