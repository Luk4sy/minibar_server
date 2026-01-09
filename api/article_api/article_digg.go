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

func (ArticleApi) ArticleDiggView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)

	var article models.ArticleModel
	err := global.DB.Take(&article, "status = ? and id = ?", enum.ArticleStatusPublished, cr.ID).Error
	if err == nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	// 判断该文章之前是否点赞过
	claims := jwts.GetClaims(c)
	var userDiggArticle models.ArticleDiggModel
	err = global.DB.Take(&userDiggArticle, "user_id = ? and article_id = ?", claims.UserID, article.ID).Error
	if err != nil {
		// 点赞
		err = global.DB.Create(&models.ArticleDiggModel{
			UserID:    claims.UserID,
			ArticleID: cr.ID,
		}).Error
		if err != nil {
			res.FailWithMsg("点赞失败", c)
			return
		}

		// TODO: 更新到缓存当中
		res.OkWithMsg("点赞成功", c)
		return
	}

	// 取消点赞
	global.DB.Model(models.ArticleDiggModel{}).Delete("user_id = ? and article_id = ?", claims.UserID, article.ID)
	res.OkWithMsg("取消点赞成功", c)
	return
}
