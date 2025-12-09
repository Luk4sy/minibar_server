package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models/enum"
	"blogx_server/utils/jwts"
	"blogx_server/utils/pwd"
	"github.com/gin-gonic/gin"
)

type UpdatePasswordRequest struct {
	OldPwd string `json:"oldPwd" binding:"required"`
	Pwd    string `json:"pwd" binding:"required"`
}

func (UserApi) UpdatePasswordView(c *gin.Context) {
	var cr UpdatePasswordRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
	}

	claims := jwts.GetClaims(c)
	user, err := claims.GetUser()
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	// 邮箱注册 or 绑定邮箱
	if !(user.RegisterSource == enum.RegisterEmailSourceType || user.Email != "") {
		res.FailWithMsg("修改失败，使用邮箱注册或绑定邮箱后课进行修改！", c)
		return
	}

	// 校验之前的密码
	if !pwd.CompareHashAndPassword(user.Password, cr.OldPwd) {
		res.FailWithMsg("旧密码错误", c)
		return
	}

	hashPwd, _ := pwd.GenerateFromPassword(cr.Pwd)
	global.DB.Model(&user).Update("password", hashPwd)
	res.OkWithMsg("密码修改成功", c)
}
