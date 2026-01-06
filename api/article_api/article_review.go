package article_api

import (
	"github.com/gin-gonic/gin"
	"minibar_server/common/res"
	"minibar_server/global"
	"minibar_server/middleware"
	"minibar_server/models"
	"minibar_server/models/enum"
)

type ArticleReviewRequest struct {
	ArticleID uint               `json:"articleID" binding:"required"`
	Status    enum.ArticleStatus `json:"status" binding:"required,oneof=3 4"`
	Msg       string             `json:"msg"` // 状态为 4 的时候，传消息给用户，审核失败，失败原因
}

func (ArticleApi) ArticleReviewView(c *gin.Context) {
	cr := middleware.GetBind[ArticleReviewRequest](c)

	var article models.ArticleModel
	err := global.DB.Take(&article, cr.ArticleID).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	global.DB.Model(&article).Update("status", cr.Status)

	// TODO: 给文章的作责发布一个系统通知

	res.FailWithMsg("审核成功", c)
}
