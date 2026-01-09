package article_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"minibar_server/common/res"
	"minibar_server/global"
	"minibar_server/middleware"
	"minibar_server/models"
)

func (ArticleApi) ArticleRemoveAdminView(c *gin.Context) {
	cr := middleware.GetBind[models.RemoveRequest](c)

	var list []models.ArticleModel
	global.DB.Find(&list, "id in ?", cr.IDList)

	if len(list) > 0 {
		err := global.DB.Delete(&list).Error
		if err != nil {
			res.FailWithMsg("删除文章失败", c)
			return
		}
	}

	res.OkWithMsg(fmt.Sprintf("删除文章成功，共计删除 %d 条", len(list)), c)
}
