package log_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"fmt"
	"github.com/gin-gonic/gin"
)

type LogApi struct {
}

type LogListRequest struct {
	Limit       int               `form:"limit"`
	Page        int               `form:"page"`
	Key         string            `form:"key"`
	LogType     enum.LogType      `form:"logType"`
	Level       enum.LogLevelType `form:"level"`
	UserID      uint              `form:"userID"`
	IP          string            `form:"ip"`
	LoginStatus bool              `form:"loginStatus"`
	ServiceName string            `form:"serviceName"`
}

type LogListResponse struct {
	models.LogModel
	UserNickname string `json:"userNickname"`
	UserAvatar   string `json:"userAvatar"`
}

func (LogApi) LogListView(c *gin.Context) {
	// 分页 查询（精确查询，模糊匹配）
	var cr LogListRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithError(err, c)
	}

	var list []models.LogModel
	if cr.Page > 20 {
		cr.Page = 1
	}
	if cr.Page <= 0 {
		cr.Page = 1
	}
	if cr.Limit == 0 || cr.Limit > 100 {
		cr.Limit = 10
	}
	offset := (cr.Page - 1) * cr.Limit

	model := models.LogModel{
		LogType:     cr.LogType,
		Level:       cr.Level,
		UserID:      cr.UserID,
		IP:          cr.IP,
		LoginStatus: cr.LoginStatus,
		ServiceName: cr.ServiceName,
	}
	like := global.DB.Where("title like ?", fmt.Sprintf("%%%s%%", cr.Key))

	global.DB.Preload("UserModel").Debug().Where(like).Where(model).Offset(offset).Limit(cr.Limit).Find(&list)

	var count int64
	global.DB.Debug().Where(like).Where(model).Model(models.LogModel{}).Count(&count)

	var _list = make([]LogListResponse, 0)
	for _, logModel := range list {
		_list = append(_list, LogListResponse{
			LogModel:     logModel,
			UserNickname: logModel.UserModel.Nickname,
			UserAvatar:   logModel.UserModel.Avatar,
		})
	}

	res.OkWithList(_list, int(count), c)
	return
}
