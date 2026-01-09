package article_api

import (
	"github.com/gin-gonic/gin"
	"minibar_server/common/res"
	"minibar_server/global"
	"minibar_server/middleware"
	"minibar_server/models"
	"minibar_server/utils/jwts"
)

func (ArticleApi) ArticleRemoveUserView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)

	claims := jwts.GetClaims(c)
	var model models.ArticleModel
	err := global.DB.Take(&model, "user_id = ? and id = ?", claims.UserID, cr.ID).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	err = global.DB.Delete(&model).Error
	if err != nil {
		res.FailWithMsg("删除文章失败", c)
		return
	}

	res.OkWithMsg("删除文章成功", c)
}
