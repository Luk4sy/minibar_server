package user_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/utils/jwts"
	"github.com/gin-gonic/gin"
	"time"
)

type UserLoginListRequest struct {
	common.PageInfo
	UserID    uint   `form:"userId"`
	Ip        string `form:"ip"`
	Addr      string `form:"addr"`
	StartTime int64  `form:"startTime"`
	EndTime   int64  `form:"endTime"`
	Type      int8   `json:"type" form:"type" binding:"required,oneof=1 2"` // 1 用户：超看自己的内容 2 管理员可以查全部
}

type UserLoginListResponse struct {
	models.UserLoginModel
	UserNickName string `json:"userNickName"`
	UserAvatar   string `json:"userAvatar"`
}

func (UserApi) UserLoginListView(c *gin.Context) {
	var cr UserLoginListRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	claim := jwts.GetClaims(c)
	if cr.Type == 1 {
		cr.UserID = claim.UserID
	}

	var query = global.DB.Where("")
	if cr.StartTime != 0 {
		t := time.Unix(cr.StartTime, 0)
		query.Where("created_at > ?", t.Format("2006-01-02 15:04:05"))
	}
	if cr.EndTime != 0 {
		t := time.Unix(cr.EndTime, 0)
		query.Where("created_at < ?", t.Format("2006-01-02 15:04:05"))
	}
	var preloads []string
	if cr.Type == 2 {
		preloads = []string{"UserModel"}
	}

	_list, count, _ := common.ListQuery(models.UserLoginModel{
		UserID: cr.UserID,
		IP:     cr.Ip,
		Addr:   cr.Addr,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Where:    query,
		Preloads: preloads,
	})

	var list = make([]UserLoginListResponse, 0)
	for _, model := range _list {
		list = append(list, UserLoginListResponse{
			UserLoginModel: model,
			UserNickName:   model.UserModel.Nickname,
			UserAvatar:     model.UserModel.Avatar,
		})
	}

	res.OkWithList(list, count, c)

}
