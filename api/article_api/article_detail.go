package article_api

import (
	"github.com/gin-gonic/gin"
	"minibar_server/common/res"
	"minibar_server/global"
	"minibar_server/middleware"
	"minibar_server/models"
	"minibar_server/models/enum"
	"minibar_server/utils/jwts"
)

type ArticleDetailResponse struct {
	models.ArticleModel
	Nickname   string `json:"nickname"`
	UserAvatar string `json:"userAvatar"`
}

func (ArticleApi) ArticleDetailView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)

	// 未登录的用户，只能看到发布成功的文章
	// 登录用户，能看到自己的所有文章
	// 管理员，能看到自己的所有文章

	var article models.ArticleModel
	err := global.DB.Preload("UserModel").Take(&article, cr.ID).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	claims, err := jwts.ParseTokenByGin(c)
	if err != nil {
		// 未登录
		if article.Status != enum.ArticleStatusPublished {
			res.FailWithMsg("文章不存在", c)
			return
		}
	}
	switch claims.Role {
	case enum.UserRole:
		if claims.UserID == article.UserID {
			// 用户查看其他人的文章
			if article.Status != enum.ArticleStatusPublished {
				res.FailWithMsg("文章不存在", c)
				return
			}
		}
	}
	// TODO:从缓存中获取浏览量和点赞数
	res.OkWithData(ArticleDetailResponse{
		ArticleModel: article,
		Nickname:     article.UserModel.Nickname,
		UserAvatar:   article.UserModel.Avatar,
	}, c)
}
